package player

import (
	"context"
	"log"
	"time"

	"fmt"

	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/leaderboard"
	"github.com/tumypmyp/chess/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Do(update tgbotapi.Update, db memory.Memory, bot Sender, cmd string) ([]Response, error) {
	player := getPlayerByID(update.SentFrom().ID, update.SentFrom().UserName, db)
	log.Println("player:", player)
	log.Println("message:", update.Message)
	if update.Message != nil && update.Message.IsCommand() {
		return player.Cmd(db, bot, update.Message)
	}
	r, err := player.Do(db, bot, cmd)
	return []Response{r}, err
}

func getPlayerByID(id int64, username string, db memory.Memory) (player Player) {
	ID := PlayerID(id)
	var err error
	if player, err = GetPlayer(ID, db); err != nil {
		player = NewPlayer(db, ID, username)
	}
	return
}

func GetPlayer(ID PlayerID, m memory.Memory) (p Player, err error) {
	key := fmt.Sprintf("user:%d", ID)
	if err = m.Get(key, &p); err != nil {
		return p, fmt.Errorf("can not get player by id: %w", err)
	}
	return
}

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
