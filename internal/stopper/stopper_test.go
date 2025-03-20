package Stopper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockStopper struct {
	allow bool
}

func (m *mockStopper) IsAllow(ip string) bool {
	return m.allow
}

func TestNewStopper(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{"Leaky Bucket", "leaky_bucket", "*Stopper.LeakyBucket"},
		{"Fixed Window", "fixed_window", "*Stopper.FixedWindow"},
		{"Sliding Window", "sliding_window", "*Stopper.SlidingWindow"},
		{"Default to Token Bucket", "unknown_mode", "*Stopper.TokenBucket"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopper := NewStopper(tt.mode, 10, time.Second)
			actual := fmt.Sprintf("%T", stopper)

			if actual != tt.expected {
				t.Errorf("NewStopper(%s) = %s, expected %s", tt.mode, actual, tt.expected)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name           string
		allowRequests  bool
		expectedStatus int
	}{
		{"Allow request", true, http.StatusOK},
		{"Block request", false, http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopper := &mockStopper{allow: tt.allowRequests}
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := Check(stopper, nextHandler)

			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
