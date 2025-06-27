package state

import (
	"fmt"
	"strings"
	"sync"
)

// это единственная штука, которая меняется
type Session struct {
	mu      sync.RWMutex
	stateMu sync.RWMutex
	state   State
	Data    sessionData
	UserID  int64
	FlowID  string
}

func NewSession(userID int64, flowID string, initialState State) *Session {
	return &Session{UserID: userID, state: initialState, Data: make(map[string]any), FlowID: flowID}
}

func (s *Session) State() State {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.state
}

func (s *Session) SetState(State State) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.state = State
}

func (s *Session) SetData(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
}

func (s *Session) GetData(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.Data[key]
	return val, ok
}

type sessionData map[string]any

func (d sessionData) String() string {
	var b strings.Builder
	for k, v := range d {
		b.WriteString(fmt.Sprintf("%s %v\n", k, v))
	}
	return b.String()
}
