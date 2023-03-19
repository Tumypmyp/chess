package helpers

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestStubBot(t *testing.T) {
	t.Run("read/write", func(t *testing.T) {
		bot := NewStubBot()
		val, err := bot.Send(tgbotapi.MessageConfig{Text: "abcd"})
		AssertNoError(t, err)
		AssertString(t, val.Text, "ok")
		AssertString(t, bot.Read(), "abcd")
	})
	t.Run("multiple read/write", func(t *testing.T) {
		bot := NewStubBot()
		val, err := bot.Send(tgbotapi.MessageConfig{Text: "abcd"})
		AssertNoError(t, err)
		AssertString(t, val.Text, "ok")

		val, err = bot.Send(tgbotapi.MessageConfig{Text: "efgh"})
		AssertNoError(t, err)
		AssertString(t, val.Text, "ok")
		AssertString(t, bot.Read(), "abcdefgh")
		AssertInt(t, bot.Len(), 2)
	})
}
