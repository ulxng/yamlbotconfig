package main

import (
	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	flowGroup := a.bot.Group()
	flowGroup.Use(a.FindFSM())

	flowGroup.Handle("/start", a.defaultStartHandler, a.startFlow)
	flowGroup.Handle(tele.OnText, a.handleError, a.handleTextFlow)
	flowGroup.Handle(tele.OnContact, a.handleError, a.handleContactFlow)
	a.bot.Handle(tele.OnCallback, a.handleButton)
}
