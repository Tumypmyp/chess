package main

import (
	"context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/player"
	pl "github.com/tumypmyp/chess/player_service/internal/player"
	"google.golang.org/grpc"
)

var db memory.Memory

type MyPlayServer struct {
	pb.UnimplementedPlayServer
}

func (p *MyPlayServer) MakePlayer(ctx context.Context, req *pb.PlayerRequest) (*empty.Empty, error) {
	pl.MakePlayer(helpers.PlayerID(req.GetPlayer().GetID()), req.GetUsername(), db)
	log.Println("making server in player")
	return &empty.Empty{}, nil
}

func (p *MyPlayServer) NewMessage(ctx context.Context, m *pb.Message) (*pb.Response, error) {
	r, err := pl.NewMessage(helpers.PlayerID(m.GetPlayer().GetID()), m.GetChatID(), m.GetCommand(), m.GetText(), db)
	log.Println("player:", r, err)
	log.Println("player:", r.Keyboard)
	// resp := &pb.Response{
	// 	Text:     r.Text,
	// 	Keyboard: r.Keyboard,
	// 	ChatsID:  r.ChatsID,
	// }
	log.Println("player sends response", r)
	return &r, nil
}

// func toKeyboard(k [][]helpers.Button) (keyboard []*pb.ArrayButton) {
// 	if k == nil {
// 		return []*pb.ArrayButton{}
// 	}
// 	keyboard = make([]*pb.ArrayButton, len(k))

// 	log.Println(keyboard)
// 	for i, v := range k {
		
// 		log.Println(v)
// 		keyboard[i] = &pb.ArrayButton{Buttons: make([]*pb.Button, len(v))}
// 		for j, b := range v {
// 			keyboard[i].Buttons[j] = &pb.Button{Text: b.Text, CallbackData: b.CallbackData}
// 		}
// 	}
// 	return 
// }

func main() {
	var err error
	db, err = memory.NewDatabase()
	
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	log.Print("...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("can't listen on 8080: %v", err)
	}
	log.Print("started linstening")

	s := grpc.NewServer()
	pb.RegisterPlayServer(s, &MyPlayServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
