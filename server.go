package main

import (
	"context"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ctx = context.Background()

var games Memory

var gameID int64

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
		if update.Message.Text == "/new_game" {
			gameID++
			games.Set("gameID", gameID)
		}
		reply(update.Message, bot)
	}

}

func sendStatus(game *Game, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(game.ChatID, game.String())

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func reply(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	log.Println(gameID)
	key := strconv.FormatInt(gameID, 10)
	var game Game
	if err := games.Get(key, &game); err != nil {
		game = Game{ChatID: message.Chat.ID}
	}

	log.Printf("game %v", game)
	game.Move(message.Text)

	if err := games.Set(key, game); err != nil {
		log.Printf("% v, could not set game", err)
	}
	sendStatus(&game, bot)

}
