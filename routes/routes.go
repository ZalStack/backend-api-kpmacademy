package routes

import (
	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/handlers"
	"backend-api-kpmacademy/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	api := app.Group("/api")
	auth := middleware.AuthProtected(cfg)
	admin := middleware.AdminOnly()

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true, "message": "KPM Academy API is running",
		})
	})

	api.Get("/testimonials", handlers.GetPublicTestimonials())
	api.Get("/packages", handlers.GetAllPackagesPublic())
	api.Get("/packages/:id", handlers.GetPackageByID())
	api.Get("/videos", handlers.GetAllVideos())
	api.Get("/videos/:id", handlers.GetVideoByID())
	api.Post("/kontak-form", handlers.CreateContactForm())
	api.Post("/support/submit", handlers.SubmitSupport())
	api.Post("/chat/send", handlers.SendChatMessage())
	api.Get("/chat/history", handlers.GetChatHistory())
	api.Post("/payment/notification", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	authGroup := api.Group("/auth")
	authGroup.Post("/register", handlers.Register(cfg))
	authGroup.Post("/login", handlers.Login(cfg))
	authGroup.Post("/refresh-token", handlers.RefreshToken(cfg))
	authGroup.Post("/forgot-password", handlers.ForgotPassword(cfg))
	authGroup.Post("/reset-password", handlers.ResetPassword(cfg))
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
