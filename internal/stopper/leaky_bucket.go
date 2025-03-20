package Stopper

import (
	"container/list"
	"sync"
	"time"
)

type LeakyBucket struct {
	mu        sync.Mutex
	queue     *list.List
	limit     int
	interval  time.Duration
	lastCheck time.Time
}

func NewLeakyBucket(limit int, interval time.Duration) *LeakyBucket {
	lk := &LeakyBucket{
		queue:     list.New(),
		limit:     limit,
		interval:  interval,
		lastCheck: time.Now(),
	}

	go lk.drain()

	return lk
}

func (lk *LeakyBucket) IsAllow(ip string) bool {
	lk.mu.Lock()
	defer lk.mu.Unlock()

	if lk.queue.Len() >= lk.limit {
		return false
	}

	lk.queue.PushBack(ip)

	return true
}

func (lk *LeakyBucket) drain() {
	for {
		time.Sleep(lk.interval)
		lk.mu.Lock()

		if lk.queue.Len() > 0 {
			lk.queue.Remove(lk.queue.Front())
		}

		lk.mu.Unlock()
	}
}
