package flow

import (
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/state"
)

// это просто носитель конфигурации. Иммутабельно
// сущность не связана с пользователем, не меняется в зависимости от стейта
type Step struct {
	State     state.State          `yaml:"state"` // todo возможно не нужно
	NextState *state.State         `yaml:"next"`
	Message   configurator.Message `yaml:"message"`
	DataCode  string               `yaml:"code"` // key для сохранения данных
	Callback  state.Callback
}
