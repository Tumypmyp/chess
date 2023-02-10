package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var games Memory

var gameID int64 = 50

var ctx = context.Background()
var player Player

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

	games, err = NewDatabase()
	if err != nil {
		panic(err)
	}

	err = games.Get("gameID", gameID)
	if err != nil {
		log.Printf("could not restore, gameID = %v", gameID)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		player.ChatID = update.Message.Chat.ID
		if update.Message.Text == "/new_game" {
			err := games.Get("gameID", gameID)
			if err != nil {
				log.Printf("could not restore, gameID = %v", gameID)
			}
			player.NewGame()
			games.Set("gameID", gameID)
		}

		reply(update.Message, bot)
	}

}

func reply(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	if player.CurrentGame() == nil {
		player.NewGame()
	}
	game := player.CurrentGame()
	log.Printf("player %+v", player)
	game.Move(message.Text)
	log.Printf("moved")
	if err := games.Set(game.ID, game); err != nil {
		log.Printf("% v, could not set game", err)
	}

	player.SendStatus(bot)

}
