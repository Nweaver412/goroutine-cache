package store

import "sync"

type Store interface {
	Set(key string, value string)
	Get(key string) string
	Delete(key string)
}

type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]string),
	}
}

func (s *InMemoryStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.store[key]
	return value, ok
}

func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
