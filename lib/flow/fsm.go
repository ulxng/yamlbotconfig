package flow

import (
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

type FSM struct {
	flow              Flow
	InitConditionFunc func(c tele.Context) bool
}

func NewFSM(flow Flow) *FSM {
	flow.InitialState = state.Initial // todo все флоу будут иметь одинаковый initial state
	return &FSM{flow: flow}
}

func (f *FSM) HandleStep(session *state.Session, input any) (Step, error) {
	// сначала обработать данные последнего стейта
	step := f.GetCurrentStep(session) // в сессии лежит последний стейт, а не будущий
	if step.DataCode != "" {
		//todo возможно добавить валиадацию ввода
		session.SetData(step.DataCode, input)
	}
	if step.NextState == nil {
		return step, ErrorEmptyNextStep
	}
	// потом переключить стейт. Пользователю будет отправляться сообщение из этого шага
	nextState := *step.NextState
	session.SetState(nextState)
	s := f.GetCurrentStep(session)
	return s, nil
}

func (f *FSM) GetCurrentStep(session *state.Session) Step {
	return f.flow.Steps[session.State()]
}

func (f *FSM) Supports(session *state.Session) bool {
	return session.FlowID == f.flow.ID
}

func (f *FSM) Start(userID int64) *state.Session {
	return state.NewSession(userID, f.flow.ID, f.flow.InitialState)
}

func (f *FSM) IsFinished(session *state.Session) bool {
	return session.State() == state.Complete
}
