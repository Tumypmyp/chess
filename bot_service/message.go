package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pb "github.com/tumypmyp/chess/proto/player"
)

func makeMessage(chatID int64, r pb.Response) (msg tgbotapi.MessageConfig) {
	msg = tgbotapi.NewMessage(chatID, r.Text)
	if len(r.Keyboard) > 0 {
		msg.ReplyMarkup = makeKeyboard(r.Keyboard)
	}
	return
}

// make inline keyboard for game
func makeKeyboard(keyboard []*pb.ArrayButton) tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, len(keyboard))

	for i, v := range keyboard {
		markup[i] = make([]tgbotapi.InlineKeyboardButton, len(v.GetButtons()))
		for j, b := range v.GetButtons() {
			markup[i][j] = tgbotapi.NewInlineKeyboardButtonData(b.GetText(), b.GetCallbackData())
		}
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}
