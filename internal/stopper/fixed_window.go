package Stopper

import (
	"sync"
	"time"
)

type FixedWindow struct {
	mu        sync.Mutex
	requests  map[string]int
	limit     int
	interval  time.Duration
	lastReset time.Time
}

func NewFixedWindow(limit int, interval time.Duration) *FixedWindow {
	return &FixedWindow{
		requests:  make(map[string]int),
		limit:     limit,
		interval:  interval,
		lastReset: time.Now(),
	}
}

func (fw *FixedWindow) IsAllow(ip string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	if now.Sub(fw.lastReset) > fw.interval {
		fw.requests = make(map[string]int)
		fw.lastReset = now
	}

	if fw.requests[ip] >= fw.limit {
		return false
	}

	fw.requests[ip]++
	return true
}
