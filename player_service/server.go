package main

import (
	"context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"
	. "github.com/tumypmyp/chess/player"
	"google.golang.org/grpc"
)

var db memory.Memory

type MyPlayServer struct {
	UnimplementedPlayServer
}

func (p *MyPlayServer) MakePlayer(ctx context.Context, player *PlayerRequest) (*empty.Empty, error) {
	MakePlayer(helpers.PlayerID(player.GetPlayer().GetID()), player.GetUsername(), db)
	log.Println("making server in player")
	return &empty.Empty{}, nil
}
func (p *MyPlayServer) NewMessage(ctx context.Context, m *Message) (*Response, error) {

	r, err := NewMessage(helpers.PlayerID(m.GetPlayer().GetID()), m.GetChatID(), m.GetCommand(), m.GetText(), db)
	log.Println(r, err)
	log.Println(toKeyboard(r.Keyboard))
	resp := &Response{
		Text:     r.Text,
		Keyboard: toKeyboard(r.Keyboard),
		ChatsID:  r.ChatsID,
	}
	log.Println(resp)
	return resp, nil
}

func toKeyboard(k [][]helpers.Button) (keyboard []*ArrayButton) {
	if k == nil {
		return []*ArrayButton{}
	}
	keyboard = make([]*ArrayButton, len(k))

	log.Println(keyboard)
	for i, v := range k {
		
		log.Println(v)
		keyboard[i] = &ArrayButton{Buttons: make([]*Button, len(v))}
		for j, b := range v {
			keyboard[i].Buttons[j] = &Button{Text: b.Text, CallbackData: b.CallbackData}
		}
	}
	return 
}

func main() {
	var err error
	db, err = memory.NewDatabase()
	
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	log.Print("...")
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("can't listen on 8888: %v", err)
	}
	log.Print("started linstening")

	s := grpc.NewServer()
	RegisterPlayServer(s, &MyPlayServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
