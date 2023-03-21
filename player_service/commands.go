package player

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/leaderboard"
	"github.com/tumypmyp/chess/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

func NewMessage(p PlayerID, chatID int64, cmd, text string, db memory.Memory) (r Response, err error) {
	if cmd != "" {
		return Cmd(db, cmd, text, p, chatID)
	}
	return Do(p, db, text, chatID)
}

// runs a command by player
func Cmd(db memory.Memory, cmd, text string, p PlayerID, ChatID int64) (r Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	log.Println(cmd, text)
	switch cmd {
	case newgame:
		r, err = doNewGame(db, p, text)
	case leaderboard:
		r1, err2 := getLeaderboard(p)
		r = r1
		err = err2
	default:
		err = NoSuchCommandError{cmd}
		r = Response{Text: err.Error(), ChatsID: []int64{ChatID}}
	}
	return
}

// Move player
func Do(id PlayerID, db memory.Memory, move string, chatID int64) (Response, error) {
	p, _ := getPlayer(id, db)
	game, err := p.CurrentGame(db)
	if err != nil {
		return Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	if err = game.Move(id, move); err != nil {
		return Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	if err := db.Set(fmt.Sprintf("game:%d", game.ID), game); err != nil {
		e := DatabaseStoringError{err}
		return Response{Text: e.Error(), ChatsID: []int64{chatID}}, e
	}
	return SendStatus(game), nil

}


// get or create new player
func MakePlayer(id PlayerID, username string, db memory.Memory) (player Player) {
	var err error
	if player, err = getPlayer(id, db); err != nil {
		player = NewPlayer(db, id, username)
	}
	return
}

// get player from memory
func getPlayer(ID PlayerID, m memory.Memory) (p Player, err error) {
	key := fmt.Sprintf("user:%d", ID)
	if err = m.Get(key, &p); err != nil {
		return p, fmt.Errorf("can not get player by id: %w", err)
	}
	return
}

type DatabaseStoringError struct {
	err error
}

func (d DatabaseStoringError) Error() string {
	return fmt.Sprintf("can not store in database: %v", d.err.Error())
}


type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

// get leaderboard with gRPC call
func getLeaderboard(id PlayerID) (Response, error) {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and return its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: fmt.Sprintf("%d", id)})
	if err != nil {
		return Response{Text: NoConnectionError{}.Error()}, NoConnectionError{}

	}
	return Response{Text: r.GetS()}, nil
}
