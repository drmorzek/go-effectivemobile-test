package framework

import "sync"

// / SessionStore хранит данные сессии
type SessionStore struct {
	sync.RWMutex
	Data map[string]map[string]interface{}
}

// NewSessionStore создает новый SessionStore
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Data: map[string]map[string]interface{}{},
	}
}

// Get возвращает данные сессии по ID
func (s *SessionStore) Get(id string) (map[string]interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	data, exists := s.Data[id]
	return data, exists
}

// Set устанавливает данные сессии по ID
func (s *SessionStore) Set(id string, sessionData map[string]interface{}) {
	s.Lock()
	defer s.Unlock()
	s.Data[id] = sessionData
}
