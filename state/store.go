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

func (s *Store) Delete(userID int64) {
	s.pool[userID] = nil
}

func (s *Store) Create(userID int64, session *Session) *Session {
	s.pool[userID] = session
	return session
}
