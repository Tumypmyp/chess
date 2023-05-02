package main

import (
	"context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/player"
	pl "github.com/tumypmyp/chess/player_service/internal/player"
	"google.golang.org/grpc"
)

var db memory.Memory

type MyPlayServer struct {
	pb.UnimplementedPlayServer
}


Move(context.Context, *MoveRequest) (*empty.Empty, error)
Status(context.Context, *GameID) (*GameStatus, error)
NewGame(context.Context, *NewGameRequest) (*GameID, error)

func (p *MyPlayServer) Move(ctx context.Context, req *pb.PlayerRequest) (*empty.Empty, error) {
	pl.MakePlayer(req.GetPlayer().GetID(), req.GetUsername(), db)
	log.Println("making server in player")
	return &empty.Empty{}, nil
}

func (p *MyPlayServer) NewMessage(ctx context.Context, m *pb.Message) (*pb.Response, error) {
	r, err := pl.NewMessage(m.GetPlayer().GetID(), m.GetChatID(), m.GetCommand(), m.GetText(), db)
	log.Println("player:", r, err)
	log.Println("player:", r.Keyboard)
	return &r, nil
}


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
