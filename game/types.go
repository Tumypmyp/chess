package game

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type PlayerID struct {
	// ChatID   int64
	UserID int64
}
