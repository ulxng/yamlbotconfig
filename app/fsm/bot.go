package fsm

import (
	"ulxng/blueprintbot/app/sender"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

type botAPI struct {
	bot    *tele.Bot
	sender sender.MessageSender
	ctx    tele.Context
}

func (e *botAPI) SetContext(c tele.Context) {
	e.ctx = c
}

func (e *botAPI) GetContext() any {
	return e.ctx
}

func (e *botAPI) SendMessage(chatID int64, message messages.Message) error {
	return e.sender.SendRaw(e.ctx, message)
}

func (e *botAPI) CallAction(action flow.Action) error {
	return e.bot.Trigger(action, e.ctx)
}

func (e *botAPI) SaveToContext(key string, value any) {
	e.ctx.Set(key, value)
}

func (e *botAPI) PrepareInput(t flow.StepType) any {
	switch t {
	case flow.TypeText:
		return e.ctx.Message().Text
	case flow.TypeContact:
		return e.ctx.Message().Contact
	default:
		return nil
	}
}
