package state

//это единственная штука, которая меняется

type Session struct {
	UserID int64
	State  State
	Data   map[string]string
	FlowID string
}
