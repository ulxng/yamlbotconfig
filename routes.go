package main

import (
	"strings"
	"ulxng/yamlbotconf/flow"

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
	a.bot.Handle(flow.SendMessage, func(c tele.Context) error {
		return a.handleFlow(c, nil)
	})
	a.bot.Handle(flow.CollectText, func(c tele.Context) error {
		input := c.Message().Text
		return a.handleFlow(c, input)
	})
	a.bot.Handle(flow.CollectContact, func(c tele.Context) error {
		input := c.Message().Contact
		return a.handleFlow(c, input)
	})
	a.bot.Handle("send_onboard_form", a.onOnboardFormComplete)
	a.bot.Handle("send_support_request", a.onSupportRequestCompete)
	a.bot.Handle("send_help_request", a.onHelpRequestCompete)
}

func (a *App) registerFlows() {
	greetFlow := a.flowRegistry.CreateFlow("greeting")
	greetFlow.InitConditionFunc = func(c tele.Context) bool {
		userID := c.Sender().ID
		user, err := a.userRepository.Find(userID)
		if err != nil {
			return false
		}
		return user == nil
	}

	supportFlow := a.flowRegistry.CreateFlow("support")
	supportFlow.InitConditionFunc = func(c tele.Context) bool {
		return c.Callback() != nil && c.Callback().Data == "support.request"
	}

	helpFlow := a.flowRegistry.CreateFlow("help")
	helpFlow.InitConditionFunc = func(c tele.Context) bool {
		if c.Callback() != nil {
			return false
		}
		s := strings.ToLower(c.Message().Text)
		return strings.Contains(s, "помощь") ||
			strings.Contains(s, "help") ||
			strings.Contains(s, "вопрос")
	}
}
