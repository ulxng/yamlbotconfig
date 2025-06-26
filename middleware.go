package main

import (
	"ulxng/yamlbotconf/flow"

	tele "gopkg.in/telebot.v4"
)

func (a *App) FindFSM() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Message().Chat.ID
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

			step := fsm.GetCurrentStep(session)
			c.Set("fsm", fsm)
			c.Set("session", session)
			c.Set("step", step)
			//передать управление на этот flow
			var action flow.Action
			if step.Action == "" {
				action = flow.SendMessage
			} else {
				action = step.Action
			}
			return a.bot.Trigger(action, c)
		}
	}
}
