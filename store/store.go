package store

import (
	"sync"
	"time"
)

// model of kv store
type Store interface {
	Set(key string, value string)  // set kv pair
	Get(key string) (string, bool) // get value by key
	Delete(key string)
}

// thread safe kv store with expiration
type TTLStore struct {
	mu     sync.RWMutex         // thread safety allowing many to one relationship
	data   map[string]string    // maps keys to values
	expiry map[string]time.Time // maps keys to their expiration times
}

func NewTTLStore() *TTLStore {
	// creates new TTLStore with initialized maps
	s := &TTLStore{
		data:   make(map[string]string),    // init empty data map
		expiry: make(map[string]time.Time), // init empty expiry map
	}
	// starts goroutine to clean up expired entries
	go s.cleanupLoop()
	return s
}

// set stores a kv pair with no time expiration
func (s *TTLStore) Set(key, value string) {
	s.mu.Lock()           // get write lock
	defer s.mu.Unlock()   // lock is released when function exits
	s.data[key] = value   // store the key-value pair
	delete(s.expiry, key) // remove any expiration for key
}

func (s *TTLStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// if key has expiration and if it has expired
	if exp, ok := s.expiry[key]; ok && time.Now().After(exp) {
		return "", false // Return empty value and false if expired
	}

	// get val from data map
	val, exists := s.data[key]
	return val, exists
}

func (s *TTLStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)   // remove kv pair from data map
	delete(s.expiry, key) // Remove TTL for key
}

func (s *TTLStore) SetWithTTL(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	s.expiry[key] = time.Now().Add(ttl)
}

// cleanupLoop runs in background goroutine to periodically remove expired entries
func (s *TTLStore) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Second) // create ticker that fires every second
	for range ticker.C {                      // loops every time ticker fires
		now := time.Now()
		s.mu.Lock()

		for k, exp := range s.expiry {
			if now.After(exp) { // if current time is after expiration
				delete(s.data, k)   // remove expired kv pair
				delete(s.expiry, k) // remove expired ttl
			}
		}

		s.mu.Unlock()
	}
}
