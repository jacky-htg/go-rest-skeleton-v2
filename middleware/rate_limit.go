package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type RateLimiter struct {
	mu          sync.Mutex
	rate        int       // maximum number of requests allowed
	burst       int       // maximum number of requests allowed in a burst
	tokens      int       // current number of tokens available
	lastChecked time.Time // last time tokens were checked
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(rate int, burst int) *RateLimiter {
	return &RateLimiter{
		rate:        rate,
		burst:       burst,
		tokens:      burst, // start with a full bucket
		lastChecked: time.Now(),
	}
}

// Allow checks if a request can proceed based on the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastChecked).Seconds()
	rl.lastChecked = now

	// Add tokens based on elapsed time
	rl.tokens += int(elapsed * float64(rl.rate))
	if rl.tokens > rl.burst {
		rl.tokens = rl.burst // cap tokens to burst size
	}

	if rl.tokens > 0 {
		rl.tokens-- // allow the request
		return true
	}

	return false // deny the request
}

func (m *Middleware) RateLimit(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rps, err := strconv.Atoi(os.Getenv("RATE_LIMIT_RPS"))
		if err != nil {
			m.Log.Error.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		burst, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
		if err != nil {
			m.Log.Error.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		limiter := NewRateLimiter(rps, (burst * rps))

		if !limiter.Allow() {
			m.Log.Error.Println("Too Many Requests")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next(w, r, ps)
	}
}
