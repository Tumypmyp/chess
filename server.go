package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	database := Memory{db}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ID := update.Message.Chat.ID
		text := update.Message.Text

		database.Set(update.Message.From.UserName, ID)

		var player Player
		database.GetPlayer(ID, &player)
		switch text[0] {
		case '/':
			player.Do(db, bot, text)
		default:
			if err := player.Move(database, text, bot); err != nil {
				player.Send(err.Error(), bot)
			}
		}
	}

}
