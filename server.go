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

		var player Player
		var err error
		if err = database.GetPlayer(ID, &player); err != nil {
			player = NewPlayer(db, ID, update.Message.From.UserName)
		}

		switch text[0] {
		case '/':
			err = player.Do(db, bot, text)
		default:
			err = player.Move(database, text, bot)
		}
		if err != nil {
			player.Send(err.Error(), bot)
		}
	}

}
