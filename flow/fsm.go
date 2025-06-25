package flow

import (
	"fmt"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/state"
)

type FSM struct {
	flow Flow
}

func NewFSM(loader *Loader, flowID string) *FSM {
	flow := loader.Flows[flowID]
	flow.InitialState = state.StateIdle // todo
	return &FSM{flow: flow}
}

func (f FSM) HandleStep(session *state.Session, input string) (configurator.Message, error) {
	step := f.flow.Steps[session.State]
	if step.DataCode != "" {
		session.Data[step.DataCode] = input
	}
	if step.Callback != nil {
		if err := step.Callback(session, input); err != nil {
			return configurator.Message{}, fmt.Errorf("step.Callback: %w", err)
		}
	}
	if step.NextState != nil {
		//todo удалить из стора? или передать наружу флаг, который скажет, что надо удалить его?
		session.State = *step.NextState
	}
	return step.Message, nil
}

func (f FSM) Supports(session *state.Session) bool {
	return session.FlowID == f.flow.ID
}

func (f FSM) Start(userID int64) *state.Session {
	session := &state.Session{UserID: userID, State: f.flow.InitialState, Data: make(map[string]string), FlowID: f.flow.ID}
	return session
}

//todo метод для добавление хендлеров
