package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/memory"
	pl "github.com/tumypmyp/chess/player"
	g "github.com/tumypmyp/chess/game"
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

		if update.Message != nil {
			text := update.Message.Text
			r, _ := pl.Do(update, db, bot, text)
			g.Send(update.SentFrom().ID, r.Text, r.Buttons, bot)
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}
			text := update.CallbackQuery.Data
			r, _ := pl.Do(update, db, bot, text)
			
			g.Send(update.SentFrom().ID, r.Text, r.Buttons, bot)
		}
	}
}
