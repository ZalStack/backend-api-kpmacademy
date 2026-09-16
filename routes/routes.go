package routes

import (
	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/handlers"
	"backend-api-kpmacademy/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true, "message": "KPM Academy API is running",
		})
	})

	api.Get("/testimonials", handlers.GetPublicTestimonials())
	api.Get("/kontak-form", handlers.GetContactForms())
	api.Post("/kontak-form", handlers.CreateContactForm())
	api.Get("/kontak-form/:id", handlers.GetContactFormByID())
	api.Post("/support/submit", handlers.SubmitSupport())
	api.Get("/support/tickets", handlers.GetSupportTickets())
	api.Post("/chat/send", handlers.SendChatMessage())
	api.Get("/chat/history", handlers.GetChatHistory())

	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register(cfg))
	auth.Post("/login", handlers.Login(cfg))
	auth.Post("/refresh-token", handlers.RefreshToken(cfg))
	auth.Post("/forgot-password", handlers.ForgotPassword(cfg))
	auth.Post("/reset-password", handlers.ResetPassword(cfg))

	protected := api.Group("")
	protected.Use(middleware.AuthProtected(cfg))

	protected.Post("/auth/logout", handlers.Logout(cfg))
	protected.Get("/profile", handlers.GetProfile())
	protected.Put("/profile", handlers.UpdateProfile())
	protected.Put("/profile/change-password", handlers.ChangePassword())

	protected.Get("/dashboard", handlers.UserDashboard())
	protected.Get("/notifications", handlers.GetUserNotifications())
	protected.Get("/notifications/unread-count", handlers.GetUnreadNotificationCount())
	protected.Post("/notifications/:id/read", handlers.MarkNotificationRead())
	protected.Post("/notifications/read-all", handlers.MarkAllNotificationsRead())

	packages := api.Group("/packages")
	packages.Get("", handlers.GetAllPackagesPublic())
	packages.Get("/:id", handlers.GetPackageByID())

	orders := protected.Group("/orders")
	orders.Post("", handlers.CreateOrder())
	orders.Get("", handlers.GetUserOrders())
	orders.Get("/status", handlers.PaymentStatus())
	orders.Get("/:id", handlers.GetOrderByID())
	orders.Post("/:id/pay", handlers.SimulatePayment())

	practice := protected.Group("/practice")
	practice.Post("/start", handlers.StartPractice())
	practice.Post("/submit", handlers.SubmitPractice())
	practice.Get("/history", handlers.GetPracticeHistory())
	practice.Get("/statistics", handlers.GetPracticeStatisticsUser())
	practice.Get("/:id", handlers.ShowPracticeSession())

	videos := protected.Group("/videos")
	videos.Get("", handlers.GetAllVideos())
	videos.Get("/:id", handlers.GetVideoByID())
	videos.Post("/:id/order", handlers.CreateVideoOrder())
	videos.Post("/:id/pay/:orderId", handlers.SimulateVideoPayment())

	protected.Post("/testimonials", handlers.CreateTestimonial())
	protected.Get("/testimonials/my", handlers.GetUserTestimonial())

	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())

	admin.Get("/dashboard", handlers.AdminDashboard())

	admin.Get("/users", handlers.GetAllUsers())
	admin.Get("/users/:id", handlers.GetUserByID())
	admin.Post("/users/:id/toggle-active", handlers.ToggleUserActive())
	admin.Get("/login-logs", handlers.GetLoginLogs())

	admin.Get("/packages", handlers.GetAllPackagesAdmin())
	admin.Post("/packages", handlers.CreatePackage())
	admin.Get("/packages/:id", handlers.GetPackageByID())
	admin.Put("/packages/:id", handlers.UpdatePackage())
	admin.Delete("/packages/:id", handlers.DeletePackage())
	admin.Post("/packages/:id/cards", handlers.AddCardToPackage())
	admin.Delete("/packages/:id/cards/:cardId", handlers.RemoveCardFromPackage())
	admin.Post("/packages/:id/import-pdf", handlers.ImportQuestionsPDF())

	admin.Get("/orders", handlers.AdminIndexOrders())
	admin.Get("/orders/:id", handlers.AdminShowOrder())
	admin.Post("/orders/:id/verify", handlers.AdminVerifyOrder())

	admin.Get("/transactions", handlers.GetAllTransactions())
	admin.Get("/transactions/stats", handlers.GetTransactionStats())
	admin.Get("/transactions/:id", handlers.ShowTransaction())

	admin.Get("/enroll-keys", handlers.GetEnrollKeys())
	admin.Get("/enroll-keys/:id", handlers.ShowEnrollKey())

	admin.Get("/practice-statistics", handlers.GetPracticeStatisticsAdmin())
	admin.Get("/practice-statistics/:id", handlers.ShowPracticeStatistics())

	admin.Get("/reports", handlers.AdminIndexReports())

	admin.Get("/testimonials", handlers.AdminIndexTestimonials())
	admin.Post("/testimonials/:id/approve", handlers.AdminApproveTestimonial())
	admin.Post("/testimonials/:id/toggle-active", handlers.AdminToggleTestimonialActive())
	admin.Delete("/testimonials/:id", handlers.AdminDeleteTestimonial())
	admin.Post("/testimonials/bulk-delete", handlers.AdminBulkDeleteTestimonials())

	admin.Get("/support", handlers.AdminIndexSupport())
	admin.Get("/support/:id", handlers.AdminShowSupport())
	admin.Post("/support/:id/answer", handlers.AdminAnswerSupport())
	admin.Put("/support/:id/status", handlers.AdminUpdateSupportStatus())
	admin.Delete("/support/:id", handlers.AdminDeleteSupport())

	admin.Get("/kontak-form", handlers.AdminIndexContactForms())
	admin.Get("/kontak-form/:id", handlers.AdminShowContactForm())
	admin.Post("/kontak-form/:id/reply", handlers.AdminReplyContactForm())
	admin.Put("/kontak-form/:id/status", handlers.AdminUpdateContactStatus())
	admin.Delete("/kontak-form/:id", handlers.AdminDeleteContactForm())

	admin.Get("/videos", handlers.GetAllVideos())
	admin.Post("/videos", handlers.CreateVideo())
	admin.Get("/videos/:id", handlers.GetVideoByID())
	admin.Put("/videos/:id", handlers.UpdateVideo())
	admin.Delete("/videos/:id", handlers.DeleteVideo())
	admin.Post("/videos/:id/toggle", handlers.ToggleVideoActive())

	admin.Get("/video-orders", handlers.GetVideoOrdersAdmin())
	admin.Post("/video-orders/:id/grant", handlers.AdminGrantVideoAccess())

	admin.Get("/notifications", handlers.AdminIndexNotifications())

	api.Get("/payment/notification", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
}
