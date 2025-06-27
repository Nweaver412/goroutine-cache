package store

import (
	"testing"
	"time"
)

func TestBasicSetGet(t *testing.T) {
	s := NewTTLStore()
	s.Set("a", "1")
	val, ok := s.Get("a")
	if !ok || val != "1" {
		t.Errorf("Expected '1', got '%s'", val)
	}
}

func TestTTLExpiration(t *testing.T) {
	s := NewTTLStore()
	s.SetWithTTL("x", "bye", 1*time.Second)

	_, ok := s.Get("x")
	if !ok {
		t.Error("Expected key to exist")
	}

	time.Sleep(2 * time.Second)
	_, ok = s.Get("x")
	if ok {
		t.Error("Expected key to expire")
	}
}
