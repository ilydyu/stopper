package Stopper

import (
	"net"
	"net/http"
	"time"
)

type Stopper interface {
	IsAllow(ip string) bool
}

func NewStopper(mode string, limit int, interval time.Duration) Stopper {
	switch mode {
	case "leaky_bucket":
		return NewLeakyBucket(limit, interval)
	case "fixed_window":
		return NewFixedWindow(limit, interval)
	case "sliding_window":
		return NewSlidingWindow(limit, interval)
	default:
		return NewTokenBucket(limit, interval)
	}
}

func Check(stopper Stopper, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if !stopper.IsAllow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
