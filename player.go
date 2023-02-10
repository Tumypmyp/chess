package main

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Player struct {
	currentGame *Game
	gamesID     []string
	ChatID      int64
	DB          Memory
}

func NewPlayer(db Memory, ChatID int64) Player {
	return Player{
		DB:     db,
		ChatID: ChatID,
	}
}
func (p *Player) CurrentGame() (*Game, error) {
	if p.currentGame == nil {
		return nil, errors.New("no current game")
	}
	return p.currentGame, nil
}
func (p *Player) NewGame() {
	var gameID int64
	err := p.DB.Get("gameID", &gameID)
	if err != nil {
		log.Printf("could not restore, gameID = %v", gameID)
	}
	p.currentGame = NewGame(p, gameID)
	gameID++
	p.DB.Set("gameID", gameID)
	p.gamesID = append(p.gamesID, p.currentGame.ID)
}

func (p *Player) Move(move string, bot *tgbotapi.BotAPI) error {
	game, err := p.CurrentGame()
	if err != nil {
		return err
	}
	if err = game.Move(move); err != nil {
		return err
	}
	if err := p.DB.Set(game.ID, game); err != nil {
		log.Printf("% v, could reach db", err)
		return err
	}
	p.SendStatus(bot)
	return nil
}
func (p *Player) Send(bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(p.ChatID, text)

	log.Printf("%+v\n%v\n", p, msg)
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("cant send: %v", err)
	}
}

func (p *Player) SendStatus(bot *tgbotapi.BotAPI) {
	game, err := p.CurrentGame()
	if err != nil {
		p.Send(bot, err.Error())
		return
	}
	p.Send(bot, game.String())
}
