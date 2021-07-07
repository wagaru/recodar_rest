package http

import (
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	pool  map[string]*rate.Limiter
	mu    sync.Mutex
	limit rate.Limit
	burst int
}

func NewRateLimiters(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		pool:  make(map[string]*rate.Limiter),
		mu:    sync.Mutex{},
		limit: r,
		burst: b,
	}
}

func (r *RateLimiter) AddIP(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.pool[ip] = rate.NewLimiter(r.limit, r.burst)
	return r.pool[ip]
}

func (r *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	limiter, exists := r.pool[ip]
	if !exists {
		r.mu.Unlock()
		return r.AddIP(ip)
	}
	r.mu.Unlock()
	return limiter
}
