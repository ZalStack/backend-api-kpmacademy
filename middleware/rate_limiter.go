package middleware

import (
	"os"
	"strconv"
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

func NewRateLimiter(limit, windowSeconds int) *RateLimiter {
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

func RateLimit(limit, windowSeconds int) fiber.Handler {
	limiter := NewRateLimiter(limit, windowSeconds)

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if !limiter.Allow(ip) {
			return utils.ErrorResponse(c, 429, "Too many requests, please try again later", nil)
		}
		return c.Next()
	}
}

var (
	AuthRateLimiter   *RateLimiter
	APIRateLimiter    *RateLimiter
	StrictRateLimiter *RateLimiter
)

func init() {
	authLimit := 30
	apiLimit := 100
	strictLimit := 15

	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "test" {
		authLimit = 500
		apiLimit = 1000
		strictLimit = 300
	}

	if val := os.Getenv("AUTH_RATE_LIMIT"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			authLimit = n
		}
	}
	if val := os.Getenv("STRICT_RATE_LIMIT"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			strictLimit = n
		}
	}

	AuthRateLimiter = NewRateLimiter(authLimit, 60)
	APIRateLimiter = NewRateLimiter(apiLimit, 60)
	StrictRateLimiter = NewRateLimiter(strictLimit, 60)
}

func AuthRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if !AuthRateLimiter.Allow(ip) {
			return utils.ErrorResponse(c, 429, "Too many authentication requests, please try again later", nil)
		}
		return c.Next()
	}
}

func StrictRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if !StrictRateLimiter.Allow(ip) {
			return utils.ErrorResponse(c, 429, "Rate limit exceeded for this endpoint", nil)
		}
		return c.Next()
	}
}
