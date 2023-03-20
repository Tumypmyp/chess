package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type PlayerID int64

type Button struct {
	Text string
	CallbackData string
}
type Response struct {
	Text string
	Keyboard [][]Button
	ChatID int64
}