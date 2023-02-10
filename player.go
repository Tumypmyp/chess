package main

import (
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
func (p *Player) CurrentGame() *Game {
	return p.currentGame
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

func (p *Player) Move(move string) {
	game := p.CurrentGame()
	game.Move(move)
	if err := player.DB.Set(game.ID, game); err != nil {
		log.Printf("% v, could not set game", err)
	}
}

func (p *Player) SendStatus(bot *tgbotapi.BotAPI) {
	text := p.currentGame.String()
	msg := tgbotapi.NewMessage(p.ChatID, text)

	log.Printf("%+v\n%v\n%+v", p, p.currentGame.String(), msg)
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("cant send: %v", err)
	}
}
