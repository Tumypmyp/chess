package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Player struct {
	GamesID  []string `json:"gamesID"`
	Username string   `json:"username"`
	ID       int64    `json:"ID"`
}

func NewPlayer(db Memory, ID int64, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	db.Set(p.Username, p.ID)
	db.SetPlayer(p.ID, p)
	return p
}

func (p *Player) CurrentGame(db Memory) (game Game, err error) {
	db.GetPlayer(p.ID, p)
	if len(p.GamesID) == 0 {
		return game, errors.New("no current game,\ntry: /new_game")
	}
	err = db.Get(p.GamesID[len(p.GamesID)-1], &game)
	return
}

func (p *Player) SetNewGame(gameID string) {
	p.GamesID = append(p.GamesID, gameID)
}

func (p *Player) NewGame(db Memory, bot Sender, playersID ...int64) (game Game) {
	gameID, _ := db.incr("gameID")
	playersID = append([]int64{p.ID}, playersID...)

	game = NewGame(db, strconv.FormatInt(gameID, 10), bot, playersID...)
	db.GetPlayer(p.ID, p)
	return
}

func (p *Player) Move(db Memory, move string, bot Sender) error {
	game, err := p.CurrentGame(db)
	if err != nil {
		return err
	}
	if err = game.Move(p.ID, move); err != nil {
		return err
	}
	if err := db.Set(game.ID, game); err != nil {
		return fmt.Errorf("could not reach db: %w", err)
	}
	game.SendStatus(db, bot)
	return nil
}

func (p *Player) Send(text string, bot Sender) {
	msg := tgbotapi.NewMessage(p.ID, text)

	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}
func (p *Player) Do(db Memory, bot Sender, cmd string) error {
	pref := "/new_game"
	if strings.HasPrefix(cmd, pref) {
		var other string
		var players []int64
		if _, err := fmt.Sscanf(cmd, "/new_game @%v", &other); err == nil {
			var id int64
			if err := db.Get(other, &id); err != nil {
				return err
			}
			players = []int64{id}
		}
		p.NewGame(db, bot, players...)
		return nil
	}
	return errors.New("no such command")
}
