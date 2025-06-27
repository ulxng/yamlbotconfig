package state

import "sync"

type Store interface {
	Get(userID int64) *Session
	Delete(session *Session)
	Save(userID int64, session *Session) *Session
}

type MemoryStore struct {
	mu   sync.RWMutex
	pool map[int64]*Session
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pool: make(map[int64]*Session),
	}
}

func (s *MemoryStore) Get(userID int64) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if session, ok := s.pool[userID]; ok {
		return session
	}
	return nil
}

func (s *MemoryStore) Delete(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pool, session.UserID)
}

func (s *MemoryStore) Save(userID int64, session *Session) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pool[userID] = session
	return session
}
