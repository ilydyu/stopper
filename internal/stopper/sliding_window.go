package Stopper

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	interval time.Duration
}

func NewSlidingWindow(limit int, interval time.Duration) *SlidingWindow {
	return &SlidingWindow{
		requests: make(map[string][]time.Time),
		limit:    limit,
		interval: interval,
	}
}

func (sw *SlidingWindow) IsAllow(ip string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.interval)

	var validRequests []time.Time
	for _, t := range sw.requests[ip] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}
	sw.requests[ip] = validRequests

	if len(validRequests) >= sw.limit {
		return false
	}

	sw.requests[ip] = append(sw.requests[ip], now)
	return true
}
