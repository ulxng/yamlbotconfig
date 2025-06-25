package flow

import (
	"ulxng/yamlbotconf/state"
)

type Manager interface {
	FindActiveFlow(userID int64) (*Flow, error)
	InitFlow(userID int64, flowID string) (*Flow, error)
	AddStateHandler(state string, handler state.Callback)
}

type StepHandler interface {
	//todo к сессии првязывать конечно не очень, но пока это простейшее решение
	HandleStep(session *state.Session, input string) error
}
