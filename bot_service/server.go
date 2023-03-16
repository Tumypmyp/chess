package main

import (
	"log"
	"os"

	"github.com/tumypmyp/chess/memory"
	"github.com/tumypmyp/chess/game"
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

	db, err := memory.NewDatabase()
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ID := game.PlayerID{
			ChatID:   update.Message.Chat.ID,
			ClientID: update.Message.From.ID,
		}
		text := update.Message.Text
		// IsCommand
		var player Player
		var err error
		if err = player.Get(ID, db); err != nil {
			player = NewPlayer(db, ID, update.Message.From.UserName)
		}

		player.Do(db, bot, text)
		
	}

}
