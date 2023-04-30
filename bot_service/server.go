package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pb "github.com/tumypmyp/chess/proto/player"
)

const key = "TELEGRAM_APITOKEN"

// initiates bot api
func NewBot() (bot *tgbotapi.BotAPI) {
	token, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("env variable %v is not set", key)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can not start bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return
}

// returns updates channel
func GetUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	return bot.GetUpdatesChan(updateConfig)
}


func main() {
	bot := NewBot()
	updates := GetUpdates(bot)

	for update := range updates {
		if update.SentFrom() == nil {
			continue
		}
		go processUpdate(update, bot)
	}
}

func processUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var text string
	if update.Message != nil {
		text = update.Message.Text
		// text = update.Message.CommandArguments()
	} else if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := bot.Request(callback); err != nil {
			log.Println(err)
		}
		text = update.CallbackQuery.Data
	}

	resp, _ := Do(update, text)
	for _, id := range resp.ChatsID {
		if _, err := bot.Send(makeMessage(id, resp)); err != nil {
			log.Printf("can't send: %v", err)
		}
	}
}



// calls player command function from update
func Do(update tgbotapi.Update, text string) (r pb.Response, err error) {
	id := update.SentFrom().ID

	var cmd string
	
	err = MakePlayer(id, update.SentFrom().UserName)
	if err != nil {
		log.Println(err)
	}

	if update.Message != nil && update.Message.IsCommand() {
		cmd = update.Message.Command()
		log.Println("cmd, text: ", cmd, text)
	}
	r, err = NewMessage(id, update.SentFrom().ID, cmd, text)

	if update.SentFrom().ID != update.FromChat().ID {
		r.ChatsID = append(r.ChatsID, update.FromChat().ID)
	}
	return r, err
}
