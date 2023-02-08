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

var gameId int64

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

	games = NewDataBase()

	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/new_game" {
			gameId++
		}
		reply(update.Message, bot)

	}

}

func reply(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	key := strconv.FormatInt(gameId, 10)
	var game Game
	if err := games.Get(key, &game); err != nil {
		game = Game{}
	}

	game.Move(message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, game.String())
	msg.ReplyToMessageID = message.MessageID

	if err := games.Set(key, game); err != nil {
		log.Printf("% v, could not send message", err)
	}

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}

}
