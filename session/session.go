package session

type State string

type Session struct {
	UserID int64
	State  State
	Data   map[State]string
	FlowID string
}
