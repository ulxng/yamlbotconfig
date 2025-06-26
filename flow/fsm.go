package flow

import (
	"fmt"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/state"

	tele "gopkg.in/telebot.v4"
)

type FSM struct {
	flow               Flow
	callbacks          map[state.State]state.Callback
	CheckInitCondition func(c tele.Context) bool
}

func NewFSM(loader *Loader, flowID string) *FSM {
	flow := loader.Flows[flowID]
	flow.InitialState = state.Initial // todo все флоу будут иметь одинаковый initial state
	return &FSM{flow: flow, callbacks: make(map[state.State]state.Callback)}
}

func (f *FSM) HandleStep(session *state.Session, input any) (configurator.Message, error) {
	// сначала обработать данные последнего стейта
	step := f.flow.Steps[session.State] // в сессии лежит последний стейт, а не будущий
	if step.DataCode != "" {
		//внимание! если key не задан - никакие данные на шаге сохраняться не будут
		session.Data[step.DataCode] = input
	}
	cb := f.GetStateCallback(session.State)
	if cb != nil {
		if err := cb(session, input); err != nil {
			return configurator.Message{}, fmt.Errorf("callback: %w", err)
		}
	}
	var message configurator.Message
	//потом переключить стейт. На текущем шаге должно отправляться сообщение того стейта, на который переключаемся
	//сообщение нужно отправлять от следующего шага
	if step.NextState != nil {
		message = f.flow.Steps[*step.NextState].Message
		session.State = *step.NextState
	}
	if session.State == state.Complete {
		//хак для последнего шага. Тк до него выполнение больше не дойдет, выполняем сразу
		cb := f.GetStateCallback(state.Complete)
		if cb != nil {
			if err := cb(session, input); err != nil {
				return configurator.Message{}, fmt.Errorf("complete callback: %w", err)
			}
		}
	}
	return message, nil
}

func (f *FSM) GetStateCallback(state state.State) state.Callback {
	return f.callbacks[state]
}

func (f *FSM) Supports(session *state.Session) bool {
	return session.FlowID == f.flow.ID
}

func (f *FSM) Start(userID int64) *state.Session {
	session := &state.Session{UserID: userID, State: f.flow.InitialState, Data: make(map[string]any), FlowID: f.flow.ID}
	return session
}

func (f *FSM) IsFinished(session *state.Session) bool {
	return session.State == state.Complete
}

func (f *FSM) SetStateCallback(state state.State, callback state.Callback) {
	f.callbacks[state] = callback
}
