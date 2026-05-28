package main

import (
	"log/slog"
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
func (s *Store) sweep() {
	now := time.Now()
	slog.Info("[Sweeping]")

	// short read lock to collect expired keys
	s.mu.RLock()
	var expired []string
	for k, e := range s.data {
		if e.hasTTL && now.After(e.expiresAt) {
			expired = append(expired, k)
		}
	}
	s.mu.RUnlock()

	if len(expired) == 0 {
		return
	}

	// write lock only for deletion
	s.mu.Lock()
	for _, k := range expired {
		delete(s.data, k)
	}
	s.mu.Unlock()
}
func (s *Store) startExpirySweep() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.sweep()
		}
	}()
}
