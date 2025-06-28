package fsm

import (
	"errors"
	"fmt"
	"log"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/messages"
	"ulxng/blueprintbot/lib/state"
)

type Executor interface {
	Handle(session *state.Session, fsm *flow.FSM) error
}

type BotAPI interface {
	SendMessage(chatID int64, message messages.Message) error
	CallAction(action flow.Action) error
	PrepareInput(t flow.StepType) any
	GetContext() any
	SaveToContext(key string, value any)
}

type BaseExecutor struct {
	store    state.Store
	registry *flow.Registry
	BotAPI
}

func NewBaseExecutor(store state.Store, registry *flow.Registry, api BotAPI) BaseExecutor {
	return BaseExecutor{store: store, registry: registry, BotAPI: api}
}

func (e *BaseExecutor) RunFSM(userID int64) error {
	session := e.store.Get(userID)
	var fsm *flow.FSM
	if session != nil {
		fsm = e.registry.FindUserActiveFlow(session)
	} else {
		fsm = e.registry.FindFlowToStart(e.BotAPI.GetContext())
		if fsm != nil {
			session = fsm.Start(userID)
			e.store.Save(userID, session)
		}
	}
	e.SaveToContext("session", session) // todo найти способ получше
	if fsm != nil {
		return e.Handle(session, fsm)
	}
	return ErrFlowNotFound
}

func (e *BaseExecutor) Handle(session *state.Session, fsm *flow.FSM) error {
	step := fsm.GetCurrentStep(session)
	nextStep, err := fsm.HandleStep(session, e.BotAPI.PrepareInput(step.Type))
	if err != nil {
		if errors.Is(err, flow.ErrorEmptyNextStep) {
			e.store.Delete(session) // prevent fsm loop
		}
		return fmt.Errorf("fsm.HandleStep: %w", err)
	}
	//экшн выполняется только на текущем шаге после сохранения ввода
	if step.Action != "" {
		if err := e.BotAPI.CallAction(step.Action); err != nil {
			log.Printf("step callAction: %q, error: %v", step.Action, err)
		}
	}
	//выполнить дефолтное действие - отправку сообщения
	if err := e.BotAPI.SendMessage(session.UserID, nextStep.Message); err != nil {
		log.Printf("sendMessage: %v", err)
	}
	if nextStep.Skip {
		return e.Handle(session, fsm)
	}
	if fsm.IsFinished(session) {
		if nextStep.Action != "" {
			if err := e.BotAPI.CallAction(nextStep.Action); err != nil {
				log.Printf("nextStep callAction: %q, error: %v", nextStep.Action, err)
			}
		}
		e.store.Delete(session)
	}
	return nil
}
