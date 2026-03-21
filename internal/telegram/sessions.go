package telegram

import "sync"

type sessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]string // chat_id → session_id
	cwds     map[int64]string // chat_id → working directory (for Task 5)
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[int64]string),
		cwds:     make(map[int64]string),
	}
}

func (s *sessionStore) Get(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[chatID]
}

func (s *sessionStore) Set(chatID int64, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = sessionID
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

func (s *sessionStore) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
}
