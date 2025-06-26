package flow

import "ulxng/blueprintbot/state"

// эта сущность должна быть в единственном экземпляре
// инстанциируется единожды
// это просто хранитель конфигурации
type Flow struct {
	ID           string               `yaml:"id"`
	Steps        map[state.State]Step `yaml:"steps"`
	InitialState state.State
}
