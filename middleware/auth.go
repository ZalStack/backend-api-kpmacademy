package middleware

import (
	"strings"

	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthProtected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.UnauthorizedResponse(c, "Invalid authorization header format")
		}

		tokenString := parts[1]
		claims, err := utils.ValidateAccessToken(tokenString, &cfg.JWT)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !utils.IsAdmin(c) {
			return utils.ForbiddenResponse(c, "Admin access required")
		}
		return c.Next()
	}
}

func OptionalAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims, err := utils.ValidateAccessToken(tokenString, &cfg.JWT)
		if err != nil {
			return c.Next()
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}
