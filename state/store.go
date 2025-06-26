package state

type Store struct {
	pool map[int64]*Session
}

func NewStore() *Store {
	return &Store{pool: make(map[int64]*Session)}
}

func (s *Store) Get(userID int64) *Session {
	if session, ok := s.pool[userID]; ok {
		return session
	}
	return nil
}

func (s *Store) Delete(session *Session) {
	s.pool[session.UserID] = nil
}

func (s *Store) Save(userID int64, session *Session) *Session {
	s.pool[userID] = session
	return session
}
