package main

import (
	"context"
	"log"
	"net"

	pb "github.com/tumypmyp/chess/leaderboard"
	"google.golang.org/grpc"
)

type LeaderboardServer struct {
	pb.UnimplementedLeaderboardServer
}

func (l *LeaderboardServer) GetLeaderboard(c context.Context, p *pb.Player) (*pb.List, error) {
	return &pb.List{S: "leaderboard..."}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterLeaderboardServer(s, &LeaderboardServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
