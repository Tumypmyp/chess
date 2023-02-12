package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Player struct {
	GamesID []string `json:"gamesID"`
	ID      int64    `json:"ID"`
	//bot     Sender
}

func NewPlayerWithBot(db Map, ID int64, bot Sender) Player {
	return Player{
		ID: ID,
		//	bot: bot,
	}
}

func NewPlayer(db Map, ID int64) Player {
	return Player{
		ID: ID,
	}
}
func (p Player) CurrentGame(db Map) (game Game, err error) {
	if len(p.GamesID) == 0 {
		return game, errors.New("no current game, try: /new_game")
	}
	err = db.Get(p.GamesID[len(p.GamesID)-1], &game)
	return
}

func (p *Player) SetNewGame(gameID string) {
	p.GamesID = append(p.GamesID, gameID)
}

func (p Player) NewGame(db Memory, playersID ...int64) *Game {
	gameID, err := db.incr("gameID")
	if err != nil {
		log.Printf("%v", err)
	}
	playersID = append([]int64{p.ID}, playersID...)

	fmt.Println("players: ", playersID)
	return NewGame(db, strconv.FormatInt(gameID, 10), playersID...)
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
