package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Cache-Control":          "no-store, no-cache, must-revalidate",
		"Pragma":                 "no-cache",
	}

	for header, expected := range expectedHeaders {
		got := resp.Header.Get(header)
		if got != expected {
			t.Errorf("Header %s = %v, want %v", header, got, expected)
		}
	}
}

func TestCORS_OriginAllowed(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")

	app := fiber.New()
	app.Use(CORS())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Origin", "http://localhost:3000")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Access-Control-Allow-Origin")
	if got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %v, want http://localhost:3000", got)
	}
}

func TestCORS_OriginNotAllowed(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	app := fiber.New()
	app.Use(CORS())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Origin", "http://evil.com")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Access-Control-Allow-Origin")
	if got != "" {
		t.Errorf("Access-Control-Allow-Origin should be empty for disallowed origin, got %v", got)
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	app := fiber.New()
	app.Use(CORS())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", http.NoBody)
	req.Header.Set("Origin", "http://localhost:3000")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("OPTIONS request status = %v, want %v", resp.StatusCode, http.StatusNoContent)
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(3, 60)

	ip := "192.168.1.1"

	if !limiter.Allow(ip) {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("Second request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("Third request should be allowed")
	}
	if limiter.Allow(ip) {
		t.Error("Fourth request should be denied")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(2, 60)

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	if !limiter.Allow(ip1) {
		t.Error("First IP first request should be allowed")
	}
	if !limiter.Allow(ip1) {
		t.Error("First IP second request should be allowed")
	}
	if limiter.Allow(ip1) {
		t.Error("First IP third request should be denied")
	}

	if !limiter.Allow(ip2) {
		t.Error("Second IP first request should be allowed")
	}
	if !limiter.Allow(ip2) {
		t.Error("Second IP second request should be allowed")
	}
	if limiter.Allow(ip2) {
		t.Error("Second IP third request should be denied")
	}
}

func TestRequestID_GeneratesUUID(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/test", func(c *fiber.Ctx) error {
		id := GetRequestID(c)
		if id == "" {
			t.Error("Expected request ID to be set")
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get(HeaderRequestID)
	if got == "" {
		t.Error("Expected X-Request-ID header to be set")
	}
}

func TestRequestID_PreservesExisting(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set(HeaderRequestID, "my-custom-id")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get(HeaderRequestID)
	if got != "my-custom-id" {
		t.Errorf("X-Request-ID = %v, want my-custom-id", got)
	}
}

func TestCORS_CredentialsHeader(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	app := fiber.New()
	app.Use(CORS())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("Origin", "http://localhost:3000")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Access-Control-Allow-Credentials")
	if got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %v, want true", got)
	}
}

func TestSecurityHeaders_PermissionsPolicy(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Permissions-Policy")
	expected := "camera=(), microphone=(), geolocation=()"
	if got != expected {
		t.Errorf("Permissions-Policy = %v, want %v", got, expected)
	}
}

func TestSecurityHeaders_HSTS(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	got := resp.Header.Get("Strict-Transport-Security")
	if got == "" {
		t.Error("Expected Strict-Transport-Security header to be set")
	}
}

func TestRateLimiter_AllowAfterReset(t *testing.T) {
	limiter := NewRateLimiter(2, 1)

	ip := "10.0.0.1"

	if !limiter.Allow(ip) {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("Second request should be allowed")
	}
	if limiter.Allow(ip) {
		t.Error("Third request should be denied")
	}
}
