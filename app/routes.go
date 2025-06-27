package main

import (
	"ulxng/blueprintbot/lib/flow"

	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	flowGroup := a.bot.Group()
	flowGroup.Use(a.fsmExecutor.Middleware())

	flowGroup.Handle("/start", a.handleError)
	flowGroup.Handle(tele.OnText, a.handleError)
	flowGroup.Handle(tele.OnCallback, a.handleError)

	//такие эндпоинты - без flow middleware
	a.bot.Handle(flow.SendMessage, func(c tele.Context) error {
		return a.fsmExecutor.HandleStep(c, nil)
	})
	a.bot.Handle(flow.CollectText, func(c tele.Context) error {
		input := c.Message().Text
		return a.fsmExecutor.HandleStep(c, input)
	})
	a.bot.Handle("send_help_request", a.onFinish)
}

func (a *App) registerFlows() {
	flow.RegisterFlow(a.flowRegistry, "gift", func(c tele.Context) bool {
		return true
	})
}
