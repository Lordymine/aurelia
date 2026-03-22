package telegram

import "sync"

type sessionEntry struct {
	sessionID string
	active    bool // true if session was used since last process start
}

type sessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*sessionEntry
	cwds     map[int64]string // chat_id → working directory
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[int64]*sessionEntry),
		cwds:     make(map[int64]string),
	}
}

func (s *sessionStore) Get(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return ""
	}
	return e.sessionID
}

func (s *sessionStore) GetWithState(chatID int64) (sessionID string, active bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return "", false
	}
	return e.sessionID, e.active
}

func (s *sessionStore) Set(chatID int64, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = &sessionEntry{sessionID: sessionID, active: true}
}

func (s *sessionStore) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
	delete(s.cwds, chatID)
}

func (s *sessionStore) GetCwd(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cwds[chatID]
}

func (s *sessionStore) SetCwd(chatID int64, cwd string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cwds[chatID] = cwd
}
