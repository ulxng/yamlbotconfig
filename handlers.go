package main

import (
	"fmt"
	"log"
	"ulxng/yamlbotconf/flow"

	tele "gopkg.in/telebot.v4"
)

func (a *App) handleStart() {
	a.bot.Handle("/start", func(c tele.Context) error {
		return a.sender.Send(c, "main.menu")
	})
}

func (a *App) handleButtons() {
	a.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		key := c.Callback().Data
		if key == "" {
			return nil
		}
		return a.sender.Edit(c, key)
	})
}

func (a *App) handleSendMailCommand() {
	a.bot.Handle("/send", func(c tele.Context) error {
		m, err := c.Bot().Send(c.Recipient(), "Отправляю сообщение на email...")
		if err != nil {
			return fmt.Errorf("bot.Send: %w", err)
		}
		go func() {
			if err := a.mailer.Send("Demo message", "Test"); err != nil {
				log.Printf("Failed to send email: %v", err)
				_, err := c.Bot().Edit(m, "Failed to send email")
				if err != nil {
					log.Printf("Failed to send message: %v", err)
					return
				}
			}
			_, err := c.Bot().Edit(m, "Заявка успешно отправлена")
			if err != nil {
				return
			}
		}()
		return nil
	})
}

func (a *App) handleFlows() {
	flowLoader := flow.NewLoader("flow")
	//todo нужен ли механизм автоматической инициализации флоу?
	//или оставить это на ручное управление?
	greetFlow := flow.NewFSM(flowLoader, "greeting")

	a.bot.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Message().Sender.ID
		session := a.store.Get(userID) // todo это может быть общая часть, а вот инициализация - нет. Ее в любом случае надо писать вручную
		if session == nil {
			session = greetFlow.Start(userID)
			a.store.Create(userID, session)
		}
		if !greetFlow.Supports(session) {
			return nil
		}
		input := c.Message().Text

		message, err := greetFlow.HandleStep(session, input)
		if err != nil {
			return fmt.Errorf("greetFlow.HandleStep: %w", err)
		}
		if greetFlow.IsFinished(session) {
			a.store.Delete(userID)
		}
		return a.sender.SendRaw(c, message)
	})
}
