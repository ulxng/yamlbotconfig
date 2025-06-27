package main

import (
	"fmt"
	"log"
	"ulxng/blueprintbot/app/model"
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

func (a *App) onOnboardFormComplete(c tele.Context) error {
	if err := a.sender.Send(c, "main.menu"); err != nil {
		log.Printf("send_onboard_form: c.Send: %v", err)
	}
	if c.Get("session") == nil {
		return nil
	}

	session := c.Get("session").(*state.Session)
	notification := fmt.Sprintf("%s %d\n", "userID", session.UserID)
	notification += session.Data.String()
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
		if err := a.mailer.Send(session.Data.String(), "Заявка на техподдержку"); err != nil {
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
		if err := a.mailer.Send(session.Data.String(), "Заявка на помощь"); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
	return nil
}
