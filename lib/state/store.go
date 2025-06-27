package state

import "sync"

type Store struct {
	mu   sync.RWMutex
	pool map[int64]*Session
}

func NewStore() *Store {
	return &Store{
		pool: make(map[int64]*Session),
	}
}

func (s *Store) Get(userID int64) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if session, ok := s.pool[userID]; ok {
		return session
	}
	return nil
}

func (s *Store) Delete(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pool, session.UserID)
}

func (s *Store) Save(userID int64, session *Session) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pool[userID] = session
	return session
}
