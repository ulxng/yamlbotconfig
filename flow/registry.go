package flow

import (
	"ulxng/yamlbotconf/state"

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

//todo как работать с flow
// нужно иметь в виду, что flow может работать сквозным - через разные обработчики
// то есть нельзя повесить flow на один обработчик и все

// todo алгоритм работы с flow
// на каждом обработчике сначала проверям - нет ли активной сессии
// если есть - значит существует какой-то запущенный flow. Его надо найти по flowID

//если сессия не найдена - должен быть понятный способ определять, нужно ли инициировать какой-то flow в этой ситуации
// то есть помимо supports по сессии должен быть supports, который проверяет условия запуска какого-то флоу
// желательно чтобы это условие хранилось внутри flow?
// но мы решили, что не пишем классов для отдельных flow. Тогда эта инфа будет храниться в самих объектах, а не структурах
// можно после NewFSM добавлять onInit condition. В него прокидывать все - context, session и тд
//
