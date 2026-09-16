package main

import (
	"log"

	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/middleware"
	"backend-api-kpmacademy/migrations"
	"backend-api-kpmacademy/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)

	migrations.RunMigrations()
	migrations.SeedAdmin()

	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.RateLimit(cfg.Rate.Limit, cfg.Rate.Expiration))

	routes.SetupRoutes(app, cfg)

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false, "message": message, "errors": err.Error(),
	})
}
