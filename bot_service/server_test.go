package main

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"
	// "github.com/tumypmyp/chess/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestDoCommand(t *testing.T) {
	t.Run("do newgame", func(t *testing.T) {
		// db := memory.NewStubDatabase()

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 123}, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err := Do(update, update.Message.Text)

		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 1)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Chat: &tgbotapi.Chat{ID: 123},
			Text: "11",
		}}

		r, err = Do(update, update.Message.Text)
		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 1)
	})

	t.Run("do newgame in other chat", func(t *testing.T) {

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 456}, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err := Do(update, update.Message.Text)

		t.Log(r)
		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 2)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Chat: &tgbotapi.Chat{ID: 456},
			Text: "11",
		}}

		r, err = Do(update, "11")
		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 2)
	})
	t.Run("do newgame with other", func(t *testing.T) {

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 123}, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err := Do(update, update.Message.Text)

		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 1)

		cmd = "/newgame"
		user = tgbotapi.User{ID: 456, UserName: "def"}
		update = tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 789}, Text: cmd + " @abc", Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err = Do(update, update.Message.Text)


		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 3)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Chat: &tgbotapi.Chat{ID: 789},
			Text: "11",
		}}

		r, err = Do(update, "11")
		AssertNoError(t, err)
		AssertInt(t, int64(len(r.ChatsID)), 3)

	})
}
