package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/leaderboard"
	"github.com/tumypmyp/chess/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Player struct {
	ID       game.PlayerID
	GamesID  []int64 `json:"gamesID"`
	Username string  `json:"username"`
}

func NewPlayer(db memory.Memory, ID game.PlayerID, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	db.Set(fmt.Sprintf("username:%v", p.Username), p.ID)
	p.Store(db)
	return p
}

type NoCurrentGameError struct{}

func (n NoCurrentGameError) Error() string { return "no current game,\ntry: /newgame" }

func (p *Player) CurrentGame(db memory.Memory) (game game.Game, err error) {
	p.Get(p.ID, db)
	if len(p.GamesID) == 0 {
		return game, NoCurrentGameError{}
	}
	err = db.Get(fmt.Sprintf("game:%d", p.GamesID[len(p.GamesID)-1]), &game)
	return
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}

func (p *Player) NewGame(db memory.Memory, bot game.Sender, players ...game.PlayerID) (current_game game.Game) {
	players = append([]game.PlayerID{p.ID}, players...)

	current_game = game.NewGame(db, bot, players...)

	for _, id := range players {
		var p Player
		err := p.Get(id, db)
		if err != nil {
			log.Println("no such player", id)
		}
		p.AddNewGame(current_game.ID)
		p.Store(db)
	}

	current_game.SendStatus(db, bot)
	p.Get(p.ID, db)
	return
}

// add p.Update()

func (p *Player) Move(db memory.Memory, bot game.Sender, move string) error {
	game, err := p.CurrentGame(db)
	if err != nil {
		return err
	}
	if err = game.Move(p.ID, move); err != nil {
		return err
	}
	if err := db.Set(fmt.Sprintf("game:%d", game.ID), game); err != nil {
		return fmt.Errorf("could not reach db: %w", err)
	}
	game.SendStatus(db, bot)
	return nil
}

func (p Player) Send(text string, bot game.Sender) {
	msg := tgbotapi.NewMessage(int64(p.ID), text)

	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}

func (p *Player) DoNewGame(db memory.Memory, bot game.Sender, cmd string) (err error) {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]
	var players []game.PlayerID
	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err = db.Get(key, &clientID); err != nil {
			return fmt.Errorf("cant find player @%v", p2)
		}

		id := game.PlayerID(clientID)

		// var player Player
		// if err := player.Get(id, db); err != nil {
		// 	id.ChatID = clientID
		// 	player.Get(id, db)
		// }
		players = append(players, id)
	}
	p.NewGame(db, bot, players...)
	return
}

func (p *Player) getLeaderboard(bot game.Sender) error {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: "Vova"})
	if err != nil {
		return fmt.Errorf("could not get leaderboard")

	}
	log.Printf("Greeting: %s", r.GetS())
	p.Send(r.GetS(), bot)
	return nil
}

func (p *Player) Do2(db memory.Memory, bot game.Sender, cmd string) error {
	pref := "/newgame"
	leaderboard := "/leaderboard"

	if strings.HasPrefix(cmd, pref) {
		return p.DoNewGame(db, bot, cmd)
	} else if strings.HasPrefix(cmd, leaderboard) {
		return p.getLeaderboard(bot)
	} else {
		return p.Move(db, bot, cmd)
	}
}
func (p *Player) Do(db memory.Memory, bot game.Sender, cmd string) error {
	err := p.Do2(db, bot, cmd)
	if err != nil {
		p.Send(err.Error(), bot)
	}
	return err
}

func (p *Player) Get(ID game.PlayerID, m memory.Memory) error {
	key := fmt.Sprintf("user:%d", ID)
	if err := m.Get(key, p); err != nil {
		return fmt.Errorf("can not get player by id: %w", err)
	}
	return nil
}

// Update memory.Memory with new value of a player
func (p Player) Store(m memory.Memory) error {
	key := fmt.Sprintf("user:%d", p.ID)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}
