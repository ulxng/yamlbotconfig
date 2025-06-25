package state

type State string

// todo пока что нужно придерживаться соглашения - первый статус и последний должны иметь такие названия
const (
	Initial  State = "idle"
	Complete State = "complete"
)

type Callback func(session *Session, input any) error
