package flow

import (
	"log"
	"ulxng/blueprintbot/lib/state"
)

type Registry struct {
	flows map[string]*FSM // todo управление приоритетом?
}

func NewRegistry(loader *Loader) *Registry {
	//все flow инициализируются при старте. Но вызов DefineFlowStart нужен, чтобы задать колбэк для старта. Иначе flow никогда не запустится
	flows := make(map[string]*FSM)
	for id, flow := range loader.Flows {
		flows[id] = NewFSM(flow)
	}
	return &Registry{flows: flows}
}

func (r *Registry) DefineFlowStart(flowID string, cb StartCondition) *FSM {
	fsm := r.GetByKey(flowID)
	if fsm == nil {
		log.Printf("defineFlowStart: unsupported flow: %q", flowID)
		return nil
	}
	fsm.setShouldStart(cb)
	return fsm
}

func (r *Registry) GetByKey(flowID string) *FSM {
	return r.flows[flowID]
}

func (r *Registry) FindUserActiveFlow(session *state.Session) *FSM {
	if session != nil {
		return r.GetByKey(session.FlowID)
	}
	return nil
}

func (r *Registry) FindFlowToStart(context any) *FSM {
	for _, fsm := range r.flows {
		if fsm.shouldStart(context) {
			return fsm
		}
	}
	return nil
}

// хелпер для инициализации флоу
// этот T - тип параметра, который будет передаваться в startCondition
func RegisterFlow[T any](r *Registry, flowID string, cond func(T) bool) *FSM {
	return r.DefineFlowStart(flowID, StartConditionFrom(cond))
}
