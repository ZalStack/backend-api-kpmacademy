package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	migrations.SeedAdmin(cfg)

	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
		BodyLimit:    int(cfg.Upload.MaxUploadSize),
	})

	app.Use(recover.New())
	app.Use(compress.New())
	if cfg.App.Debug {
		app.Use(logger.New())
	}
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.RateLimit(cfg.Rate.Limit, cfg.Rate.Expiration))

	routes.SetupRoutes(app, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
