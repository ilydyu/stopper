package Stopper

import (
	"testing"
	"time"
)

func TestTokenBucket_IsAllow(t *testing.T) {
	tests := []struct {
		name         string
		limit        int
		interval     time.Duration
		actions      []string
		expectedPass int
		expectedDeny int
	}{
		{
			name:         "Allow within limit",
			limit:        3,
			interval:     time.Second,
			actions:      []string{"allow", "allow", "allow", "deny", "deny"},
			expectedPass: 3,
			expectedDeny: 2,
		},
		{
			name:         "Reset after interval",
			limit:        3,
			interval:     time.Second,
			actions:      []string{"allow", "allow", "allow", "deny", "deny", "allow"},
			expectedPass: 4,
			expectedDeny: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tb := &TokenBucket{
				limit:     tt.limit,
				interval:  tt.interval,
				tokens:    make(map[string]int),
				lastCheck: make(map[string]time.Time),
			}

			ip := "192.168.1.1"
			var passed, denied int

			for _, action := range tt.actions {
				if action == "allow" {
					if tb.IsAllow(ip) {
						passed++
					} else {
						denied++
					}
				} else if action == "deny" {
					if !tb.IsAllow(ip) {
						denied++
					} else {
						passed++
					}
				}

				time.Sleep(200 * time.Millisecond)
			}

			if passed != tt.expectedPass {
				t.Errorf("expected %d successful requests, got %d", tt.expectedPass, passed)
			}
			if denied != tt.expectedDeny {
				t.Errorf("expected %d denied requests, got %d", tt.expectedDeny, denied)
			}
		})
	}
}
