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
			//передать управление на flow endpoint
			return e.bot.Trigger(flow.Default, c)
			// todo сразу же тут вызвать HandleStep. Зачем по контексту гонять сессию
		}
	}
}

func (e *Executor) HandleStep(c tele.Context) error {
	if c.Get("fsm") == nil {
		return ErrFlowNotFound
	}
	fsm := c.Get("fsm").(*flow.FSM)
	session := c.Get("session").(*state.Session)

	step := fsm.GetCurrentStep(session)
	nextStep, err := fsm.HandleStep(session, parseInput(c, step.Type))
	if err != nil {
		if !errors.Is(err, flow.ErrorEmptyNextStep) {
			return fmt.Errorf("fsm.HandleStep: %w", err)
		}
		//todo это может аффектить Trigger
		e.store.Delete(session)
	}
	//экшн выполняется только на текущем шаге после сохранения ввода
	if step.Action != "" {
		if err := e.bot.Trigger(step.Action, c); err != nil {
			log.Printf("bot.Trigger: action: %q, error: %v", step.Action, err)
		}
	}
	//выполнить дефолтное действие - отправку сообщения
	if err := e.sender.SendRaw(c, nextStep.Message); err != nil {
		log.Printf("sender.SendRaw: %v", err)
	}
	//после дефолтного действия выполнить кастомный экшн
	//только если есть сообщение
	if nextStep.Skip {
		return e.HandleStep(c)
	}
	if fsm.IsFinished(session) {
		//todo сделать так, чтобы в триггер прокидывалось последнее отправленное сообщение. Чтобы на нем можно было делать edit например
		if nextStep.Action != "" {
			if err := e.bot.Trigger(nextStep.Action, c); err != nil {
				log.Printf("bot.Trigger: action: %q, error: %v", nextStep.Action, err)
			}
		}
		e.store.Delete(session)
	}
	return nil
}

func parseInput(c tele.Context, t flow.StepType) any {
	switch t {
	case flow.TypeText:
		return c.Message().Text
	case flow.TypeContact:
		return c.Message().Contact
	default:
		return nil
	}
}
