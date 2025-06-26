package main

import (
	"fmt"
	"log"
	"ulxng/yamlbotconf/state"

	tele "gopkg.in/telebot.v4"
)

func (a *App) registerFlows() {
	greetFlow := a.flowRegistry.CreateFlow("greeting")
	greetFlow.SetStateCallback(state.Complete, a.greetFlowCompletedCallback)
	greetFlow.InitConditionFunc = func(c tele.Context) bool {
		return false
	}

}

func (a *App) greetFlowCompletedCallback(session *state.Session, input any) error {
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
}
