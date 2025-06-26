package main

import (
	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	a.bot.Use(a.FindFSM())

	a.bot.Handle("/start", a.defaultStartHandler, a.startFlow)
	a.bot.Handle(tele.OnCallback, a.handleButton)
	a.bot.Handle(tele.OnText, a.handleError, a.handleTextFlow)
	a.bot.Handle(tele.OnContact, a.handleError, a.handleContactFlow)
}
