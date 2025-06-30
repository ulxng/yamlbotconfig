package sender

import (
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

type Sender interface {
	Send(c tele.Context, key string) error
	Edit(c tele.Context, key string) error
	SendRaw(c tele.Context, message messages.Message) error
}

type RoutableSender interface {
	Sender
	Route(c tele.Context, buttonText string) error
}
