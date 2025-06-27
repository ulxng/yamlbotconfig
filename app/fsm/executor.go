package fsm

import (
	"errors"
	"fmt"
	"log"
	"ulxng/blueprintbot/app/sender"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

type Executor struct {
	store    state.Store
	sender   sender.MessageSender
	registry *flow.Registry
	bot      *tele.Bot
}

func NewExecutor(store state.Store, sender sender.MessageSender, registry *flow.Registry, bot *tele.Bot) *Executor {
	return &Executor{
		store:    store,
		sender:   sender,
		registry: registry,
		bot:      bot,
	}
}

func (e *Executor) Middleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Message().Chat.ID
			session := e.store.Get(userID)
			var fsm *flow.FSM
			if session != nil {
				fsm = e.registry.FindUserActiveFlow(session)
			} else {
				fsm = e.registry.FindFlowToStart(c) // todo rename
				if fsm != nil {
					session = fsm.Start(userID)
					e.store.Save(userID, session)
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
			return e.bot.Trigger(action, c)
		}
	}
}

func (e *Executor) HandleStep(c tele.Context, input any) error {
	if c.Get("fsm") == nil {
		return ErrFlowNotFound
	}
	fsm := c.Get("fsm").(*flow.FSM)
	session := c.Get("session").(*state.Session)
	step, err := fsm.HandleStep(session, input)
	if err != nil {
		if !errors.Is(err, flow.ErrorEmptyNextStep) {
			return fmt.Errorf("fsm.HandleStep: %w", err)
		}
		e.store.Delete(session)
	}
	if err := e.sender.SendRaw(c, step.Message); err != nil {
		return fmt.Errorf("sender.SendRaw: %w", err)
	}
	if fsm.IsFinished(session) {
		if step.Action != "" {
			if err := e.bot.Trigger(step.Action, c); err != nil {
				log.Printf("bot.Trigger: action: %q, error: %v", step.Action, err)
			}
		}
		e.store.Delete(session)
	}
	return nil
}
