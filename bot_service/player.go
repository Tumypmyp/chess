package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/leaderboard"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}
type PlayerID struct {
	ChatID   int64
	ClientID int64
}

type Player struct {
	ID       PlayerID
	GamesID  []int64 `json:"gamesID"`
	Username string  `json:"username"`
}

func NewPlayer(db Memory, ID PlayerID, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	db.Set(fmt.Sprintf("username:%v", p.Username), p.ID.ClientID)
	p.Store(db)
	return p
}

func (p *Player) CurrentGame(db Memory) (game Game, err error) {
	p.Get(p.ID, db)
	if len(p.GamesID) == 0 {
		return game, errors.New("no current game,\ntry: /newgame")
	}
	err = db.Get(fmt.Sprintf("game:%d", p.GamesID[len(p.GamesID)-1]), &game)
	return
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}

func (p *Player) NewGame(db Memory, bot Sender, players ...Player) (game Game) {

	players = append([]Player{*p}, players...)

	game = NewGame(db, bot, players...)
	p.Get(p.ID, db)
	return
}

// add p.Update()

func (p *Player) Move(db Memory, bot Sender, move string) error {
	game, err := p.CurrentGame(db)
	if err != nil {
		return err
	}
	if err = game.Move(*p, move); err != nil {
		return err
	}
	if err := db.Set(fmt.Sprintf("game:%d", game.ID), game); err != nil {
		return fmt.Errorf("could not reach db: %w", err)
	}
	game.SendStatus(db, bot)
	return nil
}

func (p Player) Send(text string, bot Sender) {
	msg := tgbotapi.NewMessage(p.ID.ChatID, text)

	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}
func (p *Player) DoNewGame(db Memory, bot Sender, cmd string) error {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]

	var players []Player
	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err := db.Get(key, &clientID); err != nil {
			// fmt.Printf("didnt find %v\n", p2)
			return fmt.Errorf("cant find player @%v", p2)
		}
		id := PlayerID{p.ID.ChatID, clientID}

		var player Player
		if err := player.Get(id, db); err != nil {
			id.ChatID = clientID
			player.Get(id, db)
		}
		players = append(players, player)
	}
	p.NewGame(db, bot, players...)
	return nil
}
func (p *Player) getLeaderboard(bot Sender) error {
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

func (p *Player) Do(db Memory, bot Sender, cmd string) error {
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
func (p *Player) Get(ID PlayerID, m Memory) error {
	key := fmt.Sprintf("chat:%duser:%d", ID.ChatID, ID.ClientID)
	if err := m.Get(key, p); err != nil {
		return fmt.Errorf("can not get player by id: %w", err)
	}
	return nil
}

// Update Memory with new value of a player
func (p Player) Store(m Memory) error {
	key := fmt.Sprintf("chat:%duser:%d", p.ID.ChatID, p.ID.ClientID)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}
