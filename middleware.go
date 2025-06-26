package main

import (
	"errors"
	"ulxng/yamlbotconf/flow"

	tele "gopkg.in/telebot.v4"
)

func (a *App) FindFSM() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Message().Sender.ID
			session := a.store.Get(userID)
			var fsm *flow.FSM
			if session != nil {
				fsm = a.flowRegistry.FindUserActiveFlow(session)
			} else {
				fsm = a.flowRegistry.FindFlowToStart(c) // todo rename
				if fsm != nil {
					session = fsm.Start(userID)
					a.store.Save(userID, session)
				}
			}
			if fsm == nil {
				return next(c)
			}

			c.Set("fsm", fsm)
			c.Set("session", session)
			//передать управление на этот flow
			return next(c)
		}
	}
}

func (a *App) startFlow(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if err := a.handleFlow(c, nil); err != nil {
			if errors.Is(err, ErrFlowNotFound) {
				return next(c)
			}
		}
		return nil
	}
}

func (a *App) handleTextFlow(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		input := c.Message().Text
		if err := a.handleFlow(c, input); err != nil {
			if errors.Is(err, ErrFlowNotFound) {
				return next(c)
			}
		}
		return nil
	}
}

func (a *App) handleContactFlow(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		input := c.Message().Contact
		if err := a.handleFlow(c, input); err != nil {
			if errors.Is(err, ErrFlowNotFound) {
				return next(c)
			}
		}
		return nil
	}
}
