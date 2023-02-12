package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ctx = context.Background()

var database Memory

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	database = Memory{db}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ID := update.Message.Chat.ID
		var player Player
		database.GetPlayer(ID, &player)

		switch update.Message.Text {
		case "/new_game":
			player.NewGame(database).SendStatus(database, bot)
		default:
			if err := player.Move(database, update.Message.Text, bot); err != nil {
				player.Send(err.Error(), bot)
			}
		}
	}

}
