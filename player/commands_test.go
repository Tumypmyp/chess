package player

import (
	"testing"

	"github.com/tumypmyp/chess/memory"
	"github.com/tumypmyp/chess/helpers"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestPlayerDo(t *testing.T) {
	t.Run("do newgame", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := helpers.NewStubBot()

		cmd := "/newgame"
		user := tgbotapi.User{ID: 123, UserName: "abc"}
		update := tgbotapi.Update{Message: &tgbotapi.Message{From: &user, Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}}

		err := Do(update, db, bot, cmd)

		AssertNoError(t, err)
		AssertInt(t, bot.Len(), 1)

		update = tgbotapi.Update{Message: &tgbotapi.Message{
			From: &user,
			Text: "11",
		}}

		err = Do(update, db, bot, cmd)
		AssertNoError(t, err)
		AssertInt(t, bot.Len(), 2)

	})
}
