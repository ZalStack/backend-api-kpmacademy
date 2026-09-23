package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupApp(handler fiber.Handler) *fiber.App {
	app := fiber.New()
	app.Get("/test", handler)
	return app
}

func TestSuccessResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return SuccessResponse(c, http.StatusOK, "Success", fiber.Map{"key": "value"})
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

func TestSuccessResponse_Created(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return SuccessResponse(c, http.StatusCreated, "Created", nil)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}

func TestErrorResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestNotFoundResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return NotFoundResponse(c, "Not found")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestBadRequestResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return BadRequestResponse(c, "Bad request", nil)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestUnauthorizedResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return UnauthorizedResponse(c, "Unauthorized")
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

func TestForbiddenResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return ForbiddenResponse(c, "Forbidden")
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

func TestInternalErrorResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return InternalErrorResponse(c, "Internal error")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestConflictResponse(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return ConflictResponse(c, "Conflict")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusConflict)
	}
}

func TestSuccessResponseWithPagination(t *testing.T) {
	app := setupApp(func(c *fiber.Ctx) error {
		return SuccessResponseWithPagination(c, "Paginated", []string{"a", "b"}, &Pagination{
			Total:       100,
			PerPage:     10,
			CurrentPage: 1,
			LastPage:    10,
		})
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

func TestSafeRequestID_NilContext(t *testing.T) {
	id := safeRequestID(nil)
	if id != "" {
		t.Errorf("safeRequestID(nil) = %v, want empty string", id)
	}
}

func TestSafeRequestID_NoLocals(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		id := safeRequestID(c)
		if id != "" {
			t.Errorf("safeRequestID() = %v, want empty string", id)
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close() //nolint:gocritic // test cleanup
}

func TestSafeRequestID_WithLocals(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("request_id", "test-id-123")
		id := safeRequestID(c)
		if id != "test-id-123" {
			t.Errorf("safeRequestID() = %v, want test-id-123", id)
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close() //nolint:gocritic // test cleanup
}
