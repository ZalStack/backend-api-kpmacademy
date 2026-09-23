package docs

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

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
		return c.Send(landingHTML)
	}
}
