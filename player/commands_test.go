package player

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestCommands(t *testing.T) {
	t.Run("do newgame", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := NewStubBot()

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 0}, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err := Do(update, db, bot, cmd)

		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 1)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Chat: &tgbotapi.Chat{ID: 0},
			Text: "11",
		}}

		r, err = Do(update, db, bot, "11")
		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 1)
	})
	t.Run("do newgame with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := NewStubBot()

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 0}, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err := Do(update, db, bot, cmd)
		
		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 1)

		cmd = "/newgame"
		user = tgbotapi.User{ID: 456, UserName: "def"}
		update = tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Chat: &tgbotapi.Chat{ID: 0}, Text: cmd + " @abc", Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		r, err = Do(update, db, bot, cmd)

		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 2)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Chat: &tgbotapi.Chat{ID: 0},
			Text: "11",
		}}

		r, err = Do(update, db, bot, "11")
		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 2)

	})
}
