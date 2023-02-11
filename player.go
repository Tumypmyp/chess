package main

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Player struct {
	currentGame *Game
	gamesID     []string
	ChatID      int64
	DB          Memory
	bot         Sender
}

func NewPlayerWithBot(db Map, ChatID int64, bot Sender) Player {
	return Player{
		DB:     Memory{db},
		ChatID: ChatID,
		bot:    bot,
	}
}

func NewPlayer(db Map, ChatID int64) Player {
	return Player{
		DB:     Memory{db},
		ChatID: ChatID,
	}
}
func (p *Player) CurrentGame() (*Game, error) {
	if p.currentGame == nil {
		return nil, errors.New("no current game, try: /new_game")
	}
	return p.currentGame, nil
}

func (p *Player) SetNewGame(game *Game) {
	p.currentGame = game
	p.gamesID = append(p.gamesID, p.currentGame.ID)
}

func (p *Player) NewGame(other ...*Player) *Game {
	gameID, err := p.DB.incr("gameID")
	if err != nil {
		log.Printf("%v", err)
	}
	other = append([]*Player{p}, other...)
	return NewGame(gameID, other...)
}

func (p *Player) Move(move string) error {
	game, err := p.CurrentGame()
	if err != nil {
		return err
	}
	if err = game.Move(p, move); err != nil {
		return err
	}
	if err := p.DB.Set(game.ID, game); err != nil {
		return fmt.Errorf("could not reach db: %w", err)
	}
	game.SendStatus()
	return nil
}

func (p *Player) Send(text string) {
	msg := tgbotapi.NewMessage(p.ChatID, text)
	if p.bot == nil {
		return
	}

	if _, err := p.bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}
