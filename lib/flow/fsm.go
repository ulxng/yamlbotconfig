package flow

import (
	"ulxng/blueprintbot/lib/state"
)

type FSM struct {
	flow        Flow
	shouldStart StartCondition
}

func (f *FSM) setShouldStart(shouldStart StartCondition) {
	f.shouldStart = shouldStart
}

func NewFSM(flow Flow) *FSM {
	return &FSM{flow: flow, shouldStart: func(input any) bool {
		return false
	}}
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

func (f *FSM) Start(userID int64) *state.Session {
	return state.NewSession(userID, f.flow.ID, f.flow.InitialState)
}

func (f *FSM) IsFinished(session *state.Session) bool {
	return f.GetCurrentStep(session).NextState == nil
}
