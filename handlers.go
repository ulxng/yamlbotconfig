package main

import (
	"fmt"
	"log"
	"ulxng/yamlbotconf/flow"
	"ulxng/yamlbotconf/state"

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
	flowGroup := a.bot.Group()
	flowGroup.Use(a.FindFSM())

	flowGroup.Handle(tele.OnContact, func(c tele.Context) error {
		input := c.Message().Contact
		return a.handleFlow(c, input)
	})

	flowGroup.Handle(tele.OnText, func(c tele.Context) error {
		input := c.Message().Text
		return a.handleFlow(c, input)
	})

	greetFlow := a.flowRegistry.CreateFlow("greeting")
	greetFlow.CheckInitCondition = func(c tele.Context) bool {
		return true
	}
	greetFlow.SetStateCallback(state.Complete, func(session *state.Session, input any) error {
		//todo сохранть user в бд
		notification := fmt.Sprintf("%s %d\n", "userID", session.UserID)
		for k, v := range session.Data {
			switch v.(type) {
			case string:
				notification += fmt.Sprintf("%s %s\n", k, v)
			case *tele.Contact:
				notification += fmt.Sprintf("%s %s\n", k, v.(*tele.Contact).PhoneNumber)
			}
		}
		go func() {
			if err := a.mailer.Send(notification, "Анкета"); err != nil {
				log.Printf("Failed to send email: %v", err)
			}
		}()
		return nil
	})

}

func (a *App) handleFlow(c tele.Context, input any) error {
	fsm := c.Get("fsm").(*flow.FSM)
	if fsm == nil {
		return nil
	}
	session := c.Get("session").(*state.Session)
	message, err := fsm.HandleStep(session, input)
	if err != nil {
		return fmt.Errorf("fsm.HandleStep: %w", err)
	}
	if fsm.IsFinished(session) {
		a.store.Delete(session)
	}
	return a.sender.SendRaw(c, message)
}
