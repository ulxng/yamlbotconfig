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
	flow.InitialState = state.Initial // todo все флоу будут иметь одинаковый initial state
	return &FSM{flow: flow}
}

// todo важный концептуальный вопрос - в какой момент переключать стейт
// какой стейт класть в сессию - последний или будущий?
func (f FSM) HandleStep(session *state.Session, input any) (configurator.Message, error) {
	// сначала обработать данные последнего стейта
	step := f.flow.Steps[session.State] // в сессии лежит последний стейт, а не будущий
	if step.DataCode != "" {
		//внимание! если key не задан - никакие данные на шаге сохраняться не будут
		session.Data[step.DataCode] = input
	}
	if step.Callback != nil {
		if err := step.Callback(session, input); err != nil {
			return configurator.Message{}, fmt.Errorf("step.Callback: %w", err)
		}
	}
	var message configurator.Message
	//потом переключить стейт. На текущем шаге должно отправляться сообщение того стейта, на который переключаемся
	if step.NextState != nil {
		message = f.flow.Steps[*step.NextState].Message
		session.State = *step.NextState
	}
	//пока не было ответа на стейт - не переходить на следующий
	//сообщение нужно отправлять от следующего шага
	return message, nil
}

func (f FSM) Supports(session *state.Session) bool {
	return session.FlowID == f.flow.ID
}

func (f FSM) Start(userID int64) *state.Session {
	session := &state.Session{UserID: userID, State: f.flow.InitialState, Data: make(map[string]any), FlowID: f.flow.ID}
	return session
}

func (f FSM) IsFinished(session *state.Session) bool {
	return session.State == state.Complete
}

//todo метод для добавление хендлеров
