package session

import "sync"

// Store manages session IDs and working directories per chat.
type Store struct {
	mu       sync.RWMutex
	sessions map[int64]*entry
	cwds     map[int64]string
}

type entry struct {
	sessionID string
	active    bool
}

// NewStore creates a new session store.
func NewStore() *Store {
	return &Store{
		sessions: make(map[int64]*entry),
		cwds:     make(map[int64]string),
	}
}

func (s *Store) Get(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return ""
	}
	return e.sessionID
}

func (s *Store) GetWithState(chatID int64) (sessionID string, active bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return "", false
	}
	return e.sessionID, e.active
}

func (s *Store) Set(chatID int64, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = &entry{sessionID: sessionID, active: true}
}

func (s *Store) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
	delete(s.cwds, chatID)
}

// DeactivateAll marks all sessions as inactive (cold). Used when the bridge
// process dies — sessions keep their IDs for resume, but Continue must not be
// used since the process that held them is gone.
func (s *Store) DeactivateAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.sessions {
		e.active = false
	}
}

func (s *Store) GetCwd(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cwds[chatID]
}

func (s *Store) SetCwd(chatID int64, cwd string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cwds[chatID] = cwd
}
