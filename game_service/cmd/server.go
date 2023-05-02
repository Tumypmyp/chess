package main

import (
	"context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/game"
	game "github.com/tumypmyp/chess/game_service/internal"
	"google.golang.org/grpc"
)

var db memory.Memory

type MyGameServer struct {
	pb.UnimplementedGameServer
}

func (p *MyGameServer) Move(ctx context.Context, req *pb.MoveRequest) (*empty.Empty, error) {
	err := game.Move(db, req.GetGameID(), req.GetPlayerID(), req.GetText())
	log.Println("making move in player")
	return nil, err
}

func (p *MyGameServer) Status(ctx context.Context, m *pb.GameID) (*pb.GameStatus, error) {
	r := game.MakeStatus(db, m.GetID())
	return &r, nil
}

func (p *MyGameServer) NewGame(ctx context.Context, m *pb.NewGameRequest) (*pb.GameID, error) {
	gameID := game.NewGame(db, m.GetPlayersID()...)
	return &pb.GameID{ID:gameID}, nil
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
	pb.RegisterGameServer(s, &MyGameServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
