package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Set("Pragma", "no-cache")
		return c.Next()
	}
}

func getAllowedOrigins() []string {
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsEnv == "" {
		return []string{
			"https://backend-api.kpmacademic.io",
			"https://kpmacademic.io",
			"https://www.kpmacademic.io",
			"http://localhost:3000",
			"http://localhost:5173",
		}
	}
	return strings.Split(originsEnv, ",")
}

func CORS() fiber.Handler {
	allowedOrigins := getAllowedOrigins()

	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Accept,Authorization,X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "86400")
		c.Set("Access-Control-Expose-Headers", "Content-Length")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}

func RequestTimeout(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		done := make(chan error, 1)
		go func() {
			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-time.After(timeout):
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"success": false,
				"message": "Request timeout",
			})
		}
	}
}
