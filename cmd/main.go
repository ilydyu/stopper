package main

import (
	Stopper "TokenBucket/internal/stopper"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "Host to run the proxy on")
	port := flag.Int("port", 8080, "Port to run the proxy on")
	target := flag.String("target", "https://google.com", "Target application address")
	limiterType := flag.String("limiter", "token_bucket", "Rate limiting algorithm type (token_bucket, leaky_bucket, fixed_window, sliding_window)")
	interval := flag.Duration("interval", 1*time.Minute, "Interval for rate limiting")
	limit := flag.Int("limit", 100, "Rate limit (number of requests per interval)")

	flag.Parse()

	fmt.Printf("Rate limiter type: %s, limit: %d, interval: %s\n", *limiterType, *limit, *interval)

	st := Stopper.NewStopper(*limiterType, *limit, *interval)
	handler := Stopper.Check(st, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *target, http.StatusSeeOther)
	}))

	fmt.Printf("Starting proxy on %s:%d with target %s\n", *host, *port, *target)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), handler))
}
