package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-api-kpmacademy/config"

	"github.com/gofiber/fiber/v2"
)

func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:            "test-secret-key-for-jwt-signing",
			RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
			Expiration:        3600,
			RefreshExpiration: 604800,
		},
	}
}

func TestAuthProtected_NoHeader(t *testing.T) {
	app := fiber.New()
	app.Use(AuthProtected(testConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestAuthProtected_InvalidFormat(t *testing.T) {
	app := fiber.New()
	app.Use(AuthProtected(testConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "InvalidFormat")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestAuthProtected_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(AuthProtected(testConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestAuthProtected_ValidToken(t *testing.T) {
	cfg := testConfig()
	app := fiber.New()
	app.Use(AuthProtected(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	token := generateTestToken(t, cfg)
	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestAdminOnly_NotAdmin(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_role", "user")
		return c.Next()
	})
	app.Use(AdminOnly())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
}

func TestAdminOnly_IsAdmin(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_role", "admin")
		return c.Next()
	})
	app.Use(AdminOnly())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestOptionalAuth_NoHeader(t *testing.T) {
	app := fiber.New()
	app.Use(OptionalAuth(testConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d (OptionalAuth should pass without header)", resp.StatusCode, http.StatusOK)
	}
}

func TestOptionalAuth_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(OptionalAuth(testConfig()))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d (OptionalAuth should pass on invalid token)", resp.StatusCode, http.StatusOK)
	}
}

func TestOptionalAuth_ValidToken(t *testing.T) {
	cfg := testConfig()
	app := fiber.New()
	app.Use(OptionalAuth(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		if uid == nil {
			t.Error("Expected user_id to be set in context")
		}
		return c.SendString("OK")
	})

	token := generateTestToken(t, cfg)
	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
