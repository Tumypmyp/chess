package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/memory"
	pl "github.com/tumypmyp/chess/player"
    "github.com/tumypmyp/chess/helpers"
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
		if update.SentFrom() == nil {
			continue
		}

		var text string
		if update.Message != nil {
			text = update.Message.Text
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}
			text = update.CallbackQuery.Data
		}

		resp, _ := pl.Do(update, db, bot, text)
		for _, id := range resp.ChatsID {
			sendResponse(id, resp, bot)
		}
	}
}

func sendResponse(chatID int64, r helpers.Response, bot helpers.Sender) {
	msg := tgbotapi.NewMessage(chatID, r.Text)
	msg.ReplyMarkup = makeKeyboard(r.Keyboard)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}


// make inline keyboard for game
func makeKeyboard(keyboard [][]helpers.Button) tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, len(keyboard))

	for i, v := range keyboard {
		markup[i] = make([]tgbotapi.InlineKeyboardButton, len(keyboard[i]))
		for j, _ := range v {
			markup[i][j] = tgbotapi.NewInlineKeyboardButtonData(keyboard[i][j].Text, keyboard[i][j].CallbackData)
		}
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}