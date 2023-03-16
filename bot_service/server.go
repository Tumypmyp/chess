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
		if update.Message != nil {

			ID := game.PlayerID{
				ChatID:   update.Message.Chat.ID,
				ClientID: update.Message.From.ID,
			}
			
			log.Println("textid", ID)
			text := update.Message.Text
			// IsCommand
			var player Player
			var err error
			if err = player.Get(ID, db); err != nil {
				player = NewPlayer(db, ID, update.Message.From.UserName)
			}
			log.Println("player:", player)

			player.Do(db, bot, text)
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			log.Println("message:", update.CallbackQuery.Message)
			// log.Println(update.CallbackQuery.Message.Chat.ID)
			// log.Println(update.CallbackQuery.Message.From.ID)
			ID := game.PlayerID{
				ChatID:   update.CallbackQuery.Message.Chat.ID,
				ClientID: update.SentFrom().ID,
			}
			log.Println("buttonid", ID)
			text := update.CallbackQuery.Data
			
			log.Println(text)
			var player Player
			var err error
			if err = player.Get(ID, db); err != nil {
				player = NewPlayer(db, ID, update.SentFrom().UserName)
			}
			log.Println("player:", player)
			log.Println("sending text:", text)
			player.Do(db, bot, text)
			
			log.Println("sended")
			// And finally, send a message containing the data received.
			
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			
			log.Println("sent new message")
		}
	}

}
