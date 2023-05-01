package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/tumypmyp/chess/proto/player"
)

func MakePlayer(id int64, username string) error {
	conn, err := grpc.Dial("player:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPlayClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = c.MakePlayer(ctx, &pb.PlayerRequest{Username: username, Player: &pb.PlayerID{ID: int64(id)}})
	
	if err != nil {
		log.Println(err)
	}
	return err
}

func NewMessage(id int64, chatID int64, cmd, text string) (r pb.Response, err error) {
	conn, err := grpc.Dial("player:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPlayClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := c.NewMessage(ctx, &pb.Message{
		Player:  &pb.PlayerID{ID: int64(id)},
		ChatID:  chatID,
		Command: cmd,
		Text:    text,
	})

	if err != nil {
		log.Println(err)
	}
	log.Printf("server got: %v\n", res)
	return *res, err
}
