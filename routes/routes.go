package routes

import (
	"time"

	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/docs"
	"backend-api-kpmacademy/handlers"
	"backend-api-kpmacademy/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Root welcome endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Welcome to KPM Academy API",
			"version": "v1.0.0",
			"docs": fiber.Map{
				"landing":   "/docs",
				"scalar":    "/docs/scalar",
				"reference": "/docs/reference",
				"swagger":   "/docs/swagger",
				"redoc":     "/docs/redoc",
				"yaml":      "/docs/swagger.yaml",
				"json":      "/docs/openapi.json",
			},
			"api_v1": "/api/v1",
		})
	})

	// Modern API Documentation Endpoints
	app.Get("/docs", docs.LandingPageHandler())
	app.Get("/docs/scalar", docs.ScalarUIHandler())
	app.Get("/docs/reference", docs.ScalarUIHandler())
	app.Get("/docs/swagger", docs.SwaggerUIHandler())
	app.Get("/docs/redoc", docs.RedocUIHandler())
	app.Get("/docs/swagger.yaml", docs.SwaggerYAMLHandler())
	app.Get("/docs/openapi.json", docs.OpenAPIJSONHandler())

	// Primary versioned API routes (/api/v1)
	v1 := app.Group("/api/v1")
	registerAPIRoutes(v1, cfg)

	// Backward compatibility alias (/api)
	legacy := app.Group("/api")
	registerAPIRoutes(legacy, cfg)

	// Fallback 404 handler for undefined routes
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
		})
	})
}

func registerAPIRoutes(api fiber.Router, cfg *config.Config) {
	auth := middleware.AuthProtected(cfg)
	admin := middleware.AdminOnly()

	// API info endpoint
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "KPM Academy API v1",
			"endpoints": fiber.Map{
				"health":       "/health",
				"auth":         "/auth",
				"profile":      "/profile",
				"dashboard":    "/dashboard",
				"packages":     "/packages",
				"videos":       "/videos",
				"orders":       "/orders",
				"practice":     "/practice",
				"testimonials": "/testimonials",
				"kontak-form":  "/kontak-form",
				"support":      "/support",
				"chat":         "/chat",
				"admin":        "/admin",
				"docs":         "/docs",
			},
		})
	})

	api.Get("/health", func(c *fiber.Ctx) error {
		db := database.GetDB()
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Database connection error",
			})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Database unreachable",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "KPM Academy API is running",
			"version": "v1",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	api.Get("/testimonials", handlers.GetPublicTestimonials())
	api.Get("/packages", handlers.GetAllPackagesPublic())
	api.Get("/packages/:id", handlers.GetPackageByID())
	api.Get("/videos", handlers.GetAllVideos())
	api.Get("/videos/:id", handlers.GetVideoByID())
	api.Post("/kontak-form", middleware.AuthRateLimit(), handlers.CreateContactForm())
	api.Post("/support/submit", middleware.AuthRateLimit(), handlers.SubmitSupport())
	api.Post("/chat/send", middleware.AuthRateLimit(), handlers.SendChatMessage())
	api.Get("/chat/history", auth, handlers.GetChatHistory())
	api.Post("/payment/notification", handlers.MidtransWebhook())
	api.Post("/payment/video-notification", handlers.MidtransVideoWebhook())

	authGroup := api.Group("/auth")
	authGroup.Post("/register", middleware.AuthRateLimit(), handlers.Register(cfg))
	authGroup.Post("/login", middleware.AuthRateLimit(), handlers.Login(cfg))
	authGroup.Post("/refresh-token", handlers.RefreshToken(cfg))
	authGroup.Post("/forgot-password", middleware.StrictRateLimit(), handlers.ForgotPassword(cfg))
	authGroup.Post("/reset-password", middleware.StrictRateLimit(), handlers.ResetPassword(cfg))
	authGroup.Post("/logout", auth, handlers.Logout(cfg))

	api.Get("/profile", auth, handlers.GetProfile())
	api.Put("/profile", auth, handlers.UpdateProfile())
	api.Put("/profile/change-password", auth, handlers.ChangePassword())

	api.Get("/dashboard", auth, handlers.UserDashboard())

	api.Get("/notifications", auth, handlers.GetUserNotifications())
	api.Get("/notifications/unread-count", auth, handlers.GetUnreadNotificationCount())
	api.Post("/notifications/:id/read", auth, handlers.MarkNotificationRead())
	api.Post("/notifications/read-all", auth, handlers.MarkAllNotificationsRead())

	api.Get("/kontak-form", auth, handlers.GetContactForms())
	api.Get("/kontak-form/:id", auth, handlers.GetContactFormByID())

	api.Post("/orders", auth, handlers.CreateOrder())
	api.Get("/orders", auth, handlers.GetUserOrders())
	api.Get("/orders/status", auth, handlers.PaymentStatus())
	api.Get("/orders/:id", auth, handlers.GetOrderByID())
	api.Post("/orders/:id/pay", auth, handlers.SimulatePayment())

	api.Post("/practice/start", auth, handlers.StartPractice())
	api.Post("/practice/submit", auth, handlers.SubmitPractice())
	api.Get("/practice/history", auth, handlers.GetPracticeHistory())
	api.Get("/practice/statistics", auth, handlers.GetPracticeStatisticsUser())
	api.Get("/practice/:id", auth, handlers.ShowPracticeSession())

	api.Post("/videos/:id/order", auth, handlers.CreateVideoOrder())
	api.Post("/videos/:id/pay/:orderId", auth, handlers.SimulateVideoPayment())

	api.Post("/testimonials", auth, handlers.CreateTestimonial())
	api.Get("/testimonials/my", auth, handlers.GetUserTestimonial())

	adminRoutes := api.Group("/admin", auth, admin)

	adminRoutes.Get("/dashboard", handlers.AdminDashboard())

	adminRoutes.Get("/users", handlers.GetAllUsers())
	adminRoutes.Get("/users/:id", handlers.GetUserByID())
	adminRoutes.Post("/users/:id/toggle-active", handlers.ToggleUserActive())
	adminRoutes.Get("/login-logs", handlers.GetLoginLogs())

	adminRoutes.Get("/packages", handlers.GetAllPackagesAdmin())
	adminRoutes.Post("/packages", handlers.CreatePackage())
	adminRoutes.Get("/packages/:id", handlers.GetPackageByID())
	adminRoutes.Put("/packages/:id", handlers.UpdatePackage())
	adminRoutes.Delete("/packages/:id", handlers.DeletePackage())
	adminRoutes.Post("/packages/:id/cards", handlers.AddCardToPackage())
	adminRoutes.Delete("/packages/:id/cards/:cardId", handlers.RemoveCardFromPackage())
	adminRoutes.Post("/packages/:id/import-pdf", handlers.ImportQuestionsPDF())

	adminRoutes.Get("/orders", handlers.AdminIndexOrders())
	adminRoutes.Get("/orders/:id", handlers.AdminShowOrder())
	adminRoutes.Post("/orders/:id/verify", handlers.AdminVerifyOrder())

	adminRoutes.Get("/transactions", handlers.GetAllTransactions())
	adminRoutes.Get("/transactions/stats", handlers.GetTransactionStats())
	adminRoutes.Get("/transactions/:id", handlers.ShowTransaction())

	adminRoutes.Get("/enroll-keys", handlers.GetEnrollKeys())
	adminRoutes.Get("/enroll-keys/:id", handlers.ShowEnrollKey())

	adminRoutes.Get("/practice-statistics", handlers.GetPracticeStatisticsAdmin())
	adminRoutes.Get("/practice-statistics/:id", handlers.ShowPracticeStatistics())

	adminRoutes.Get("/reports", handlers.AdminIndexReports())

	adminRoutes.Get("/testimonials", handlers.AdminIndexTestimonials())
	adminRoutes.Post("/testimonials/:id/approve", handlers.AdminApproveTestimonial())
	adminRoutes.Post("/testimonials/:id/toggle-active", handlers.AdminToggleTestimonialActive())
	adminRoutes.Delete("/testimonials/:id", handlers.AdminDeleteTestimonial())
	adminRoutes.Post("/testimonials/bulk-delete", handlers.AdminBulkDeleteTestimonials())

	adminRoutes.Get("/support", handlers.AdminIndexSupport())
	adminRoutes.Get("/support/:id", handlers.AdminShowSupport())
	adminRoutes.Post("/support/:id/answer", handlers.AdminAnswerSupport())
	adminRoutes.Put("/support/:id/status", handlers.AdminUpdateSupportStatus())
	adminRoutes.Delete("/support/:id", handlers.AdminDeleteSupport())

	adminRoutes.Get("/kontak-form", handlers.AdminIndexContactForms())
	adminRoutes.Get("/kontak-form/:id", handlers.AdminShowContactForm())
	adminRoutes.Post("/kontak-form/:id/reply", handlers.AdminReplyContactForm())
	adminRoutes.Put("/kontak-form/:id/status", handlers.AdminUpdateContactStatus())
	adminRoutes.Delete("/kontak-form/:id", handlers.AdminDeleteContactForm())

	adminRoutes.Get("/videos", handlers.GetAllVideos())
	adminRoutes.Post("/videos", handlers.CreateVideo())
	adminRoutes.Get("/videos/:id", handlers.GetVideoByID())
	adminRoutes.Put("/videos/:id", handlers.UpdateVideo())
	adminRoutes.Delete("/videos/:id", handlers.DeleteVideo())
	adminRoutes.Post("/videos/:id/toggle", handlers.ToggleVideoActive())

	adminRoutes.Get("/video-orders", handlers.GetVideoOrdersAdmin())
	adminRoutes.Post("/video-orders/:id/grant", handlers.AdminGrantVideoAccess())

	adminRoutes.Get("/notifications", handlers.AdminIndexNotifications())
}
