package state

type State string

const StateIdle State = "idle"

type Callback func(session *Session, input string) error
