package main

import (
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

func (a *App) handleError(c tele.Context) error {
	return a.sender.Send(c, "errors.unknown_action")
}

func (a *App) onFinish(c tele.Context) error {
	if c.Get("session") == nil {
		return nil
	}
	session := c.Get("session").(*state.Session)
	return c.Send(session.Data.String())
}
