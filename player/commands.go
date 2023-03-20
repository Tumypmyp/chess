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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Do(update tgbotapi.Update, db memory.Memory, bot Sender, cmd string) ([]Response, error) {
	player := makePlayer(update.SentFrom().ID, update.SentFrom().UserName, db)
	log.Println("player:", player)
	log.Println("message:", update.Message)
	if update.Message != nil && update.Message.IsCommand() {
		return Cmd(db, update.Message, player, update.FromChat().ID)
	}
	r, err := player.Do(db, cmd, update.FromChat().ID)
	log.Println(r, err,cmd)
	return r, err
}



type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

// runs a command by player
func  Cmd(db memory.Memory, cmd *tgbotapi.Message, p Player, ChatID int64) (r []Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	switch cmd.Command() {
	case newgame:
		r, err = doNewGame(db, p, cmd.Text)
	case leaderboard:
		r1, err2 := getLeaderboard(p)
		r = []Response{r1}
		err = err2
	default:
		err = NoSuchCommandError{cmd.Command()}
		r = []Response{{Text: err.Error(), ChatID : ChatID}}
	}
	return
}

// get or create new player
func makePlayer(id int64, username string, db memory.Memory) (player Player) {
	ID := PlayerID(id)
	var err error
	if player, err = getPlayer(ID, db); err != nil {
		player = NewPlayer(db, ID, username)
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


type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

// get leaderboard with gRPC call
func getLeaderboard(p Player) (Response, error) {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: fmt.Sprintf("%d", p.ID)})
	if err != nil {
		return Response{Text: NoConnectionError{}.Error()}, NoConnectionError{}

	}
	log.Printf("Greeting: %s", r.GetS())
	return Response{Text: r.GetS()}, nil
}
