package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"
	pl "github.com/tumypmyp/chess/player"
)

// initiates bot api
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

		resp, _ := Do(update, db, text)
		for _, id := range resp.ChatsID {
			sendResponse(id, resp, bot)
		}
	}
}

// calls player command function from update
func Do(update tgbotapi.Update, db memory.Memory, cmd string) (r helpers.Response, err error) {
	id := helpers.PlayerID(update.SentFrom().ID)
	log.Println("player:", id)
	log.Println("message:", update.Message)

	pl.MakePlayer(id, update.SentFrom().UserName, db)
	if update.Message != nil && update.Message.IsCommand() {
		r, err = pl.Cmd(db, update.Message.Command(), update.Message.Text, id, update.SentFrom().ID)
	} else {
		r, err = pl.Do(id, db, cmd, update.SentFrom().ID)
	}
	if update.SentFrom().ID != update.FromChat().ID {
		r.ChatsID = append(r.ChatsID, update.FromChat().ID)
	}
	log.Println(cmd, r, err)
	return r, err
}

// sends response message to bot api
func sendResponse(chatID int64, r helpers.Response, bot Sender) {
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
