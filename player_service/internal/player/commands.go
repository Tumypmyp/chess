package player

import (
	"context"
	"fmt"
	"time"

	. "github.com/tumypmyp/chess/helpers"
	
	g "github.com/tumypmyp/chess/player_service/internal/game"
	"github.com/tumypmyp/chess/proto/leaderboard"
	pb "github.com/tumypmyp/chess/proto/player"
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

func NewMessage(p PlayerID, chatID int64, cmd, text string, db memory.Memory) (r pb.Response, err error) {
	if cmd != "" {
		return Cmd(db, cmd, text, p, chatID)
	}
	return Do(p, db, text, chatID)
}

// runs a command by player
func Cmd(db memory.Memory, cmd, text string, p PlayerID, ChatID int64) (r pb.Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	// log.Println(cmd, text)
	switch cmd {
	case newgame:
		r, err = doNewGame(db, p, text)
	case leaderboard:
		r, err = getLeaderboard(p, ChatID)
	default:
		err = NoSuchCommandError{cmd}
		r = pb.Response{Text: err.Error(), ChatsID: []int64{ChatID}}
	}
	return
}


type DatabaseStoringError struct {
	err error
}

func (d DatabaseStoringError) Error() string {
	return fmt.Sprintf("can not store in database: %v", d.err.Error())
}


// Move player
func Do(id PlayerID, db memory.Memory, move string, chatID int64) (pb.Response, error) {
	// p, _ := getPlayer(id, db)
	game, err := CurrentGame(id, db)
	if err != nil {
		return pb.Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	if err = game.Move(id, move); err != nil {
		return pb.Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	if err := db.Set(fmt.Sprintf("game:%d", game.ID), game); err != nil {
		e := DatabaseStoringError{err}
		return pb.Response{Text: e.Error(), ChatsID: []int64{chatID}}, e
	}
	return g.SendStatus(game), nil

}


// get or create new player
func MakePlayer(id PlayerID, username string, db memory.Memory) (player Player) {
	var err error
	if player, err = getPlayer(id, db); err != nil {
		player = NewPlayer(db, id, username)
	}
	return
}




type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

// get leaderboard with gRPC call
func getLeaderboard(id PlayerID, ChatID int64) (pb.Response, error) {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return pb.Response{Text: NoConnectionError{}.Error()}, NoConnectionError{}
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and return its pb.Response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: fmt.Sprintf("%d", id)})
	if err != nil {
		return pb.Response{Text: NoConnectionError{}.Error(), ChatsID: []int64{ChatID}}, NoConnectionError{}

	}
	return pb.Response{Text: r.GetS(), ChatsID: []int64{ChatID}}, nil
}
