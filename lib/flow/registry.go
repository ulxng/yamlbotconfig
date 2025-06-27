package flow

import (
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

type Registry struct {
	flows  map[string]*FSM
	loader *Loader
}

func NewRegistry(loader *Loader) *Registry {
	return &Registry{loader: loader, flows: make(map[string]*FSM)}
}

func (r *Registry) CreateFlow(flowID string) *FSM {
	f := NewFSM(r.loader, flowID)
	r.flows[flowID] = f
	return f
}

func (r *Registry) FindUserActiveFlow(session *state.Session) *FSM {
	if session != nil {
		for _, fsm := range r.flows {
			if fsm.Supports(session) {
				return fsm
			}
		}
	}
	return nil
}

func (r *Registry) FindFlowToStart(c tele.Context) *FSM {
	for _, fsm := range r.flows {
		if fsm.InitConditionFunc(c) {
			return fsm
		}
	}
	return nil
}
