package Stopper

import (
	"testing"
	"time"
)

func TestLeakyBucket_IsAllow(t *testing.T) {
	lk := NewLeakyBucket(3, 200*time.Millisecond)
	ip := "192.168.1.1"

	if !lk.IsAllow(ip) {
		t.Errorf("Expected first request to be allowed")
	}
	if !lk.IsAllow(ip) {
		t.Errorf("Expected second request to be allowed")
	}
	if !lk.IsAllow(ip) {
		t.Errorf("Expected third request to be allowed")
	}
	if lk.IsAllow(ip) {
		t.Errorf("Expected fourth request to be denied")
	}

	time.Sleep(250 * time.Millisecond)

	if !lk.IsAllow(ip) {
		t.Errorf("Expected first request after drain to be allowed")
	}
}
