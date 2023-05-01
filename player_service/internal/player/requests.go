package player

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	game "github.com/tumypmyp/chess/proto/game"
)

func makeNewGame(id ...int64) (gameID int64) {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := game.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	game, err := c.NewGame(ctx, &game.NewGameRequest{PlayersID: id})
	gameID = game.ID
	if err != nil {
		log.Println(err)
	}
	return gameID
}

func Move(gameID, playerID int64, text string) error {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := game.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = c.Move(ctx, &game.MoveRequest{GameID: gameID, PlayerID: playerID, Text: text})

	if err != nil {
		log.Println(err)
	}
	return err
}



func makeStatus(gameID int64) game.GameStatus {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := game.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	status, err := c.Status(ctx, &game.GameID{ID: gameID})

	if err != nil {
		log.Println(err)
	}
	return *status
}