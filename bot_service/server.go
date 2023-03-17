package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/memory"
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
		ID := game.PlayerID{
			ChatID:   update.FromChat().ID,
			ClientID: update.SentFrom().ID,
		}
		// IsCommand
		var player Player
		var err error
		if err = player.Get(ID, db); err != nil {
			player = NewPlayer(db, ID, update.Message.From.UserName)
		}
		log.Println("player:", player)

		if update.Message != nil {
			text := update.Message.Text
			player.Do(db, bot, text)
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}
			text := update.CallbackQuery.Data
			player.Do(db, bot, text)
		}
	}
}
