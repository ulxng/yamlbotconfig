package main

import (
	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	flowGroup := a.bot.Group()
	flowGroup.Use(a.FindFSM())

	flowGroup.Handle("/start", a.menuHandler)
	flowGroup.Handle(tele.OnCallback, a.handleButton)
	flowGroup.Handle(tele.OnContact, a.handleError)
	flowGroup.Handle(tele.OnText, a.handleError)

	//такие эндпоинты - без flow middleware
	a.bot.Handle("send_message", func(c tele.Context) error {
		return a.handleFlow(c, nil)
	})
	a.bot.Handle("collect_answer", func(c tele.Context) error {
		input := c.Message().Text
		return a.handleFlow(c, input)
	})
	a.bot.Handle("handle_phone", func(c tele.Context) error {
		input := c.Message().Contact
		return a.handleFlow(c, input)
	})
	a.bot.Handle("send_support_request", func(c tele.Context) error {
		return c.Send("blabla")
	})
}
