package docs

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// defaultAPIBaseURL is the production API base URL
const defaultAPIBaseURL = "https://backend-api.kpmacademic.io/api/v1"

// getAPIBaseURL returns the configured API base URL from APP_URL env.
// If APP_URL already contains /api/v1 it is used as-is, otherwise /api/v1 is appended.
// This allows the docs to stay in sync with deployment without rebuilding.
func getAPIBaseURL() string {
	raw := strings.TrimSpace(os.Getenv("APP_URL"))
	if raw == "" {
		return defaultAPIBaseURL
	}
	raw = strings.TrimSuffix(raw, "/")
	if strings.Contains(raw, "/api/v1") || strings.Contains(raw, "/api") {
		return raw
	}
	return raw + "/api/v1"
}

//go:embed swagger.yaml
var swaggerYAML []byte

//go:embed openapi.json
var openAPIJSON []byte

//go:embed swagger.html
var swaggerHTML []byte

//go:embed redoc.html
var redocHTML []byte

//go:embed scalar.html
var scalarHTML []byte

//go:embed landing.html
var landingHTML []byte

// SwaggerYAMLHandler serves the raw OpenAPI YAML file
func SwaggerYAMLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/yaml; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=3600")
		apiURL := getAPIBaseURL()
		if apiURL != defaultAPIBaseURL {
			// Dynamically replace servers URL when APP_URL differs from default
			out := bytes.ReplaceAll(swaggerYAML, []byte(defaultAPIBaseURL), []byte(apiURL))
			// Also handle legacy localhost for backwards compat if still present
			out = bytes.ReplaceAll(out, []byte("http://localhost:3000/api/v1"), []byte(apiURL))
			return c.Send(out)
		}
		return c.Send(swaggerYAML)
	}
}

// SwaggerHandler keeps backward compatibility with previous route name
func SwaggerHandler() fiber.Handler {
	return SwaggerYAMLHandler()
}

// OpenAPIJSONHandler serves the OpenAPI specification in JSON format
func OpenAPIJSONHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=3600")
		apiURL := getAPIBaseURL()
		if apiURL != defaultAPIBaseURL {
			out := bytes.ReplaceAll(openAPIJSON, []byte(defaultAPIBaseURL), []byte(apiURL))
			out = bytes.ReplaceAll(out, []byte("http://localhost:3000/api/v1"), []byte(apiURL))
			return c.Send(out)
		}
		return c.Send(openAPIJSON)
	}
}

// SwaggerUIHandler serves the classic Swagger UI
func SwaggerUIHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(swaggerHTML)
	}
}

// RedocUIHandler serves the clean publication-style Redoc UI
func RedocUIHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(redocHTML)
	}
}

// ScalarUIHandler serves the modern 3-column interactive Scalar API Reference
func ScalarUIHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(scalarHTML)
	}
}

// LandingPageHandler serves the modern documentation landing page
func LandingPageHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		apiURL := getAPIBaseURL()
		if apiURL != defaultAPIBaseURL {
			out := bytes.ReplaceAll(landingHTML, []byte(defaultAPIBaseURL), []byte(apiURL))
			out = bytes.ReplaceAll(out, []byte("http://localhost:3000/api/v1"), []byte(apiURL))
			return c.Send(out)
		}
		return c.Send(landingHTML)
	}
}
