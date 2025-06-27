package state

import (
	"sync"
)

// это единственная сущность, которая меняется
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
