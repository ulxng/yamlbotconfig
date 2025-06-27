package flow

import (
	"ulxng/blueprintbot/lib/messages"
	"ulxng/blueprintbot/lib/state"
)

// это просто носитель конфигурации. Иммутабельно
// сущность не связана с пользователем, не меняется в зависимости от стейта
type Step struct {
	NextState *state.State     `yaml:"next"`
	Message   messages.Message `yaml:"message"`
	DataCode  string           `yaml:"code"` // key для сохранения данных
	Action    Action           `yaml:"action"`
}

// эта сущность должна быть в единственном экземпляре
// инстанциируется единожды
// это просто хранитель конфигурации
type Flow struct {
	ID           string               `yaml:"id"`
	Steps        map[state.State]Step `yaml:"steps"`
	InitialState state.State          `yaml:"initial"`
}

type Action = string

// набор стандартных экшнов на шаги флоу. Кастомные тоже можно использовать, но эти зафиксированы
const (
	SendMessage    Action = "send_message"
	CollectText    Action = "collect_text"
	CollectContact Action = "collect_contact"
)

type StartCondition func(any) bool

// адаптер для унификации условия старта
// лучше не читать. добавлено, чтобы отвязать сигнатуру условия от конкретной реализации telegram bot api
func StartConditionFrom[T any](fn func(T) bool) StartCondition {
	return func(v any) bool {
		t, ok := v.(T)
		if !ok {
			return false
		}
		return fn(t)
	}
}
