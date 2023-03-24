package main

import (
	"context"
	"log"
	"net"

	"github.com/tumypmyp/chess/proto/leaderboard"
	"google.golang.org/grpc"
)

type LeaderboardServer struct {
	leaderboard.UnimplementedLeaderboardServer
}

func (l *LeaderboardServer) GetLeaderboard(c context.Context, p *leaderboard.Player) (*leaderboard.List, error) {
	return &leaderboard.List{S: "leaderboard...\n1.\n2.\n3."}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	leaderboard.RegisterLeaderboardServer(s, &LeaderboardServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
