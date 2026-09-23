package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	HeaderRequestID     = "X-Request-ID"
	ContextKeyRequestID = "request_id"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(HeaderRequestID, requestID)
		c.Locals(ContextKeyRequestID, requestID)

		return c.Next()
	}
}

func GetRequestID(c *fiber.Ctx) string {
	id, ok := c.Locals(ContextKeyRequestID).(string)
	if !ok {
		return ""
	}
	return id
}
