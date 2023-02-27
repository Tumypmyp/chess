package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}
type PlayerID struct {
	ChatID int64
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
		return game, errors.New("no current game,\ntry: /new_game")
	}
	err = db.Get(fmt.Sprintf("game:%d", p.GamesID[len(p.GamesID)-1]), &game)
	return
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}

func (p *Player) NewGame(db Memory, bot Sender, playersID ...PlayerID) (game Game) {

	playersID = append([]PlayerID{p.ID}, playersID...)

	game = NewGame(db, bot, playersID...)
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
	n, _ := fmt.Sscanf(cmd, "/new_game @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]

	var playersID []PlayerID
	for _, p2 := range others {
		var id int64
		key := fmt.Sprintf("username:%v", p2)
		if err := db.Get(key, &id); err != nil {
			// fmt.Printf("didnt find %v\n", p2)
			return fmt.Errorf("cant find player @%v", p2)
		}

		playersID = append(playersID, PlayerID{id, id})
	}
	p.NewGame(db, bot, playersID...)
	return nil
}
func (p *Player) Do(db Memory, bot Sender, cmd string) error {
	pref := "/new_game"

	if strings.HasPrefix(cmd, pref) {
		return p.DoNewGame(db, bot, cmd)
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
