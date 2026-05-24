package main

import (
	"sync"
	"time"
)

type Entry struct {
	value     string
	expiresAt time.Time
	hasTTL    bool
}

type Store struct {
	data map[string]Entry

	mu sync.RWMutex
}

func NewStore() Store {
	return Store{
		data: map[string]Entry{},
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := Entry{value: value, hasTTL: ttl != 0}
	if e.hasTTL {
		e.expiresAt = time.Now().Add(ttl)
	}
	s.data[key] = e
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return "", false
	}
	if val.hasTTL && time.Now().After(val.expiresAt) {
		return "", false // expired, treat as missing
	}
	return val.value, ok
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}
