package main

import (
	"fmt"
	"log"
	"ulxng/yamlbotconf/model"
	"ulxng/yamlbotconf/state"

	tele "gopkg.in/telebot.v4"
)

func (a *App) registerFlows() {
	greetFlow := a.flowRegistry.CreateFlow("greeting")
	//todo remove. Use action callbacks instead
	greetFlow.SetStateCallback(state.Complete, a.greetFlowCompletedCallback)
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

}

func (a *App) greetFlowCompletedCallback(session *state.Session, input any) error {
	notification := fmt.Sprintf("%s %d\n", "userID", session.UserID)
	for k, v := range session.Data {
		switch v.(type) {
		case string:
			notification += fmt.Sprintf("%s %s\n", k, v)
		case *tele.Contact:
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
