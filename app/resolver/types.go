package resolver

import (
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

// все что делает эта штука
// находит сообщение по идентификатору
// - преобразует Message в нужный формат, который понимает telebot
type MessageResolver interface {
	Get(key string) (message interface{}, markup *tele.ReplyMarkup, err error)
	Convert(msgLayout messages.Message) (message interface{}, markup *tele.ReplyMarkup, err error)
}

// в нем должна быть вся инфа про сообщения
type ReplyButtonResolver interface {
	FindNextByText(buttonText string) (message interface{}, markup *tele.ReplyMarkup, err error)
}

type RoutableResolver interface {
	MessageResolver
	ReplyButtonResolver
}
