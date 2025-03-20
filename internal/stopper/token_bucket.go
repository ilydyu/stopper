package Stopper

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu        sync.Mutex
	tokens    map[string]int
	lastCheck map[string]time.Time
	limit     int
	interval  time.Duration
}

func NewTokenBucket(limit int, interval time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:    make(map[string]int),
		lastCheck: make(map[string]time.Time),
		limit:     limit,
		interval:  interval,
	}
}

func (s *TokenBucket) IsAllow(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	last, exists := s.lastCheck[ip]

	if !exists || now.Sub(last) > s.interval {
		s.tokens[ip] = s.limit
		s.lastCheck[ip] = now
	}

	if s.tokens[ip] > 0 {
		s.tokens[ip]--
		return true
	}

	return false
}
