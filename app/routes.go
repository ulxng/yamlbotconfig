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
	// обязательный хендлер в проекте. Без него не будет работать fsm. Однако если не нужны кастомные обработчики - этого достаточно, чтобы флоу корректно работал
	a.bot.Handle(flow.Default, a.fsmExecutor.HandleStep)
	a.bot.Handle("send_help_request", a.onFinish)
}

func (a *App) registerFlows() {
	flow.RegisterFlow(a.flowRegistry, "gift", func(c tele.Context) bool {
		return true
	})
}
