package Stopper

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSlidingWindowLimiter_Concurrent(t *testing.T) {
	limiter := NewSlidingWindow(5, time.Second)
	ip := "192.168.1.1"

	var wg sync.WaitGroup
	passed := int32(0)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.IsAllow(ip) {
				atomic.AddInt32(&passed, 1)
			}
		}()
	}

	wg.Wait()

	if passed > 5 {
		t.Errorf("More requests were missed than should have been: %d", passed)
	}
}
