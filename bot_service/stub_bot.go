package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StubBot struct {
	buffer *bytes.Buffer
}

func NewStubBot() StubBot {
	return StubBot{
		buffer: new(bytes.Buffer),
	}
}

func (s StubBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	value, ok := c.(tgbotapi.MessageConfig)
	if !ok {
		return tgbotapi.Message{}, fmt.Errorf("cant convert to message")
	}
	text := value.Text
	fmt.Fprint(s.buffer, text)
	return tgbotapi.Message{Text: "ok"}, nil
}

func (s StubBot) Read() (string) {
	return s.buffer.String()
}