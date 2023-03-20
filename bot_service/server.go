package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/memory"
	pl "github.com/tumypmyp/chess/player"
    "github.com/tumypmyp/chess/helpers"
)

func NewBot() (bot *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return
}

func main() {
	bot := NewBot()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	db, err := memory.NewDatabase()
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	for update := range updates {
		if update.SentFrom() == nil {
			continue
		}

		var resp []helpers.Response
		if update.Message != nil {
			text := update.Message.Text
			resp, _ = pl.Do(update, db, bot, text)

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}
			text := update.CallbackQuery.Data
			resp, _ = pl.Do(update, db, bot, text)
		}
		for _, r := range resp {
			Send(update.SentFrom().ID, r.Text, &r.Buttons, bot)
		}
	}
}

func Send(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup, bot helpers.Sender) {
	msg := tgbotapi.NewMessage(chatID, text)
	// msg.ReplyMarkup = keyboard
	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}

