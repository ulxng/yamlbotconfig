package flow

import (
	"ulxng/yamlbotconf/state"

	tele "gopkg.in/telebot.v4"
)

type FSM struct {
	flow              Flow
	InitConditionFunc func(c tele.Context) bool
}

func NewFSM(loader *Loader, flowID string) *FSM {
	flow := loader.Flows[flowID]
	flow.InitialState = state.Initial // todo все флоу будут иметь одинаковый initial state
	return &FSM{flow: flow}
}

func (f *FSM) HandleStep(session *state.Session, input any) (*Step, error) {
	// сначала обработать данные последнего стейта
	step := f.flow.Steps[session.State] // в сессии лежит последний стейт, а не будущий
	if step.DataCode != "" {
		//внимание! если key не задан - никакие данные на шаге сохраняться не будут
		session.Data[step.DataCode] = input
	}
	var s Step
	//потом переключить стейт. На текущем шаге должно отправляться сообщение того стейта, на который переключаемся
	//сообщение нужно отправлять от следующего шага
	if step.NextState != nil {
		s = f.flow.Steps[*step.NextState]
		session.State = *step.NextState
	}
	return &s, nil
}

func (f *FSM) GetCurrentStep(session *state.Session) Step {
	return f.flow.Steps[session.State]
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
