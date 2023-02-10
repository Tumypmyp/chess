
package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Player struct {
	currentGame *Game
	gamesID     []string
	ChatID      int64
}

func (p *Player) CurrentGame() *Game {
	return p.currentGame
}
func (p *Player) NewGame() {
	p.currentGame = NewGame(p)
	p.gamesID = append(p.gamesID, p.currentGame.ID)
}

func (p *Player) SendStatus(bot *tgbotapi.BotAPI) {
	text := p.currentGame.String()
	msg := tgbotapi.NewMessage(p.ChatID, text)

	log.Printf("%+v\n%v\n%+v", p, p.currentGame.String(), msg)
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("cant send: %v", err)
	}
}
