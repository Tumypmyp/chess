package main

import (
	"context"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var games map[int64]*Game

var ctx = context.Background()

var db Memory

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

	games = make(map[int64]*Game)

	db = NewDataBase()

	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		reply(update.Message, bot)
	}

}

func reply(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	key := strconv.FormatInt(message.Chat.ID, 10)
	var game Game
	if err := db.Get(key, &game); err != nil {
		game = Game{}
	}

	game.Move(message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, game.Render())
	msg.ReplyToMessageID = message.MessageID

	if err := db.Set(key, game); err != nil {
		log.Printf("% v, could not send message", err)
	}

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}

}
