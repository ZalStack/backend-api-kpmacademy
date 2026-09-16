package middleware

import (
	"sync"
	"time"

	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
)

type Visitor struct {
	Count     int
	LastReset time.Time
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, windowSeconds int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
		window:   time.Duration(windowSeconds) * time.Second,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.LastReset) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists || time.Since(v.LastReset) > rl.window {
		rl.visitors[ip] = &Visitor{
			Count:     1,
			LastReset: time.Now(),
		}
		return true
	}

	if v.Count >= rl.limit {
		return false
	}

	v.Count++
	return true
}

func RateLimit(limit int, windowSeconds int) fiber.Handler {
	limiter := NewRateLimiter(limit, windowSeconds)

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if !limiter.Allow(ip) {
			return utils.ErrorResponse(c, 429, "Too many requests, please try again later", nil)
		}
		return c.Next()
	}
}
