package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ctx = context.Background()

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

	games, err := NewDatabase()
	if err != nil {
		panic(err)
	}

	player := NewPlayer(games, 0)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		player.ChatID = update.Message.Chat.ID

		switch update.Message.Text {
		case "/new_game":
			player.NewGame()
			player.SendStatus(bot)
		default:
			move(player, update.Message, bot)
		}
	}

}

func move(player Player, message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	err := player.Move(message.Text, bot)
	if err != nil {
		player.Send(bot, err.Error())
	}
}
