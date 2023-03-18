package player

import (
	"bytes"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StubBot struct {
	buffer *bytes.Buffer
	len    *int64
}

func NewStubBot() StubBot {
	return StubBot{
		buffer: new(bytes.Buffer),
		len:    new(int64),
	}
}

func (s StubBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	*s.len++
	value, ok := c.(tgbotapi.MessageConfig)
	if !ok {
		return tgbotapi.Message{}, fmt.Errorf("cant convert to message")
	}
	text := value.Text
	fmt.Fprint(s.buffer, text)
	return tgbotapi.Message{Text: "ok"}, nil
}

func (s StubBot) Read() string {
	return s.buffer.String()
}

func (s StubBot) Len() int64 {
	return *s.len
}
