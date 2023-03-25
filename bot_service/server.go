package main

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/helpers"
	pb "github.com/tumypmyp/chess/proto/player"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// returns updates channel
func GetUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	return bot.GetUpdatesChan(updateConfig)
}

// sends response message via bot
func SendResponse(chatID int64, r pb.Response, bot Sender) {
	msg := tgbotapi.NewMessage(chatID, r.Text)
	if len(r.Keyboard) > 0 {
		msg.ReplyMarkup = makeKeyboard(r.Keyboard)
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
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
	} else if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := bot.Request(callback); err != nil {
			log.Println(err)
		}
		text = update.CallbackQuery.Data
	}

	resp, _ := Do(update, text)
	for _, id := range resp.ChatsID {
		SendResponse(id, resp, bot)
	}
}

// calls player command function from update
func Do(update tgbotapi.Update, text string) (r pb.Response, err error) {
	id := helpers.PlayerID(update.SentFrom().ID)

	var cmd string
	log.Println("making player")
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


func NewMessage(id helpers.PlayerID, chatID int64, cmd, text string) (r pb.Response, err error) {
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c :=  pb.NewPlayClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	res, err := c.NewMessage(ctx, &pb.Message{
		Player: &pb.PlayerID{ID:int64(id)},
		ChatID: chatID,
		Command:cmd,
		Text   : text,
	})

	if err != nil {
		log.Println(err)
	}
	log.Printf("server got: %v\n", res)
	return *res, err
}

func MakePlayer(id helpers.PlayerID, username string) error {

	log.Println("connecting to 8888")
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c :=  pb.NewPlayClient(conn)

	log.Println("new client")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	log.Println("request")
	_, err = c.MakePlayer(ctx, & pb.PlayerRequest{Username: username, Player: & pb.PlayerID{ID:int64(id)}})

	log.Println("response ", err)
	return err
}



// make inline keyboard for game
func makeKeyboard(keyboard []* pb.ArrayButton) tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, len(keyboard))

	for i, v := range keyboard {
		markup[i] = make([]tgbotapi.InlineKeyboardButton, len(v.GetButtons()))
		for j, b := range v.GetButtons() {
			markup[i][j] = tgbotapi.NewInlineKeyboardButtonData(b.GetText(), b.GetCallbackData())
		}
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}
