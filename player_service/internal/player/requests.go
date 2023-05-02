package player

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/tumypmyp/chess/proto/game"
)

func makeNewGame(id ...int64) (gameID int64, err error) {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	c := game.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	g, err := c.NewGame(ctx, &game.NewGameRequest{PlayersID: id})
	if err != nil {
		return 
	}
	return g.ID, nil
}

func Move(gameID, playerID int64, text string) (err error) {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
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



func makeStatus(gameID int64) (status *game.GameStatus, err error) {
	conn, err := grpc.Dial("game:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	c := game.NewGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	status, err = c.Status(ctx, &game.GameID{ID: gameID})

	if err != nil {
		return 
	}
	return status, nil
}