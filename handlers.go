package main

import (
	"fmt"
	"ulxng/yamlbotconf/flow"
	"ulxng/yamlbotconf/state"

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
	if fsm.IsFinished(session) {
		a.store.Delete(session)
	}
	return a.sender.SendRaw(c, step.Message)
}
