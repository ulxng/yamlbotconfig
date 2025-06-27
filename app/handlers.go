package main

import (
	"fmt"
	"log"
	"ulxng/blueprintbot/app/model"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

func (a *App) menuHandler(c tele.Context) error {
	return a.sender.Send(c, "main.menu")
}

func (a *App) handleButton(c tele.Context) error {
	key := c.Callback().Data
	if key == "" {
		return nil
	}
	return a.sender.Edit(c, key)
}

func (a *App) handleError(c tele.Context) error {
	return a.sender.Send(c, "errors.unknown_action")
}

func (a *App) handleFlow(c tele.Context, input any) error {
	if c.Get("fsm") == nil {
		return ErrFlowNotFound
	}
	fsm := c.Get("fsm").(*flow.FSM)
	session := c.Get("session").(*state.Session)
	step, err := fsm.HandleStep(session, input)
	if err != nil {
		return fmt.Errorf("fsm.HandleStep: %w", err)
	}
	if err := a.sender.SendRaw(c, step.Message); err != nil {
		return fmt.Errorf("sender.SendRaw: %w", err)
	}
	if fsm.IsFinished(session) {
		if step.Action != "" {
			if err := a.bot.Trigger(step.Action, c); err != nil {
				log.Printf("bot.Trigger: action: %q, error: %v", step.Action, err)
			}
		}
		a.store.Delete(session)
	}
	return nil
}

func (a *App) onOnboardFormComplete(c tele.Context) error {
	if err := a.sender.Send(c, "main.menu"); err != nil {
		log.Printf("send_onboard_form: c.Send: %v", err)
	}
	if c.Get("session") == nil {
		return nil
	}

	session := c.Get("session").(*state.Session)
	notification := fmt.Sprintf("%s %d\n", "userID", session.UserID)
	for k, v := range session.Data {
		switch v.(type) {
		case string:
			notification += fmt.Sprintf("%s %s\n", k, v)
		case *tele.Contact:
			contact := v.(*tele.Contact)
			if contact == nil {
				continue
			}
			notification += fmt.Sprintf("%s %s\n", k, v.(*tele.Contact).PhoneNumber)
		}
	}
	if err := a.userRepository.CreateUser(model.User{ID: session.UserID}); err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	go func() {
		if err := a.mailer.Send(notification, "Анкета"); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
	return nil
}

func (a *App) onSupportRequestCompete(c tele.Context) error {
	//todo отправить сообщение админу в личный чат
	if c.Get("session") == nil {
		return nil
	}
	session := c.Get("session").(*state.Session)
	go func() {
		if err := a.mailer.Send(fmt.Sprintf("%v", session.Data), "Заявка на техподдержку"); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
	return nil
}

func (a *App) onHelpRequestCompete(c tele.Context) error {
	if c.Get("session") == nil {
		return nil
	}
	session := c.Get("session").(*state.Session)
	go func() {
		if err := a.mailer.Send(fmt.Sprintf("%v", session.Data), "Заявка на помощь"); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
	return nil
}
