package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}

		var req models.CreateOrderRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if req.PackageID == "" {
			return utils.BadRequestResponse(c, "package_id is required", nil)
		}

		db := database.GetDB()
		pkgID, err := uuid.Parse(req.PackageID)
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}

		var pkg models.Package
		if result := db.Where("id = ? AND is_active = ?", pkgID, true).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found or inactive")
		}

		totalPrice := pkg.Price
		if pkg.IsDiscountActive && pkg.DiscountPrice != nil {
			totalPrice = *pkg.DiscountPrice
		}
		if req.TotalPrice > 0 && pkg.IsPayWhatYouWant {
			if req.TotalPrice < pkg.MinPayAmount {
				return utils.BadRequestResponse(c, "Minimum payment is required", fiber.Map{
					"min_amount": pkg.MinPayAmount,
				})
			}
			totalPrice = req.TotalPrice
		}

		orderNumber := "ORD-" + time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(6)

		order := models.Order{
			UserID:         uid,
			PackageID:      &pkgID,
			Type:           "package",
			OrderNumber:    orderNumber,
			TotalPrice:     totalPrice,
			IsCustomAmount: req.IsCustomAmount,
			PaymentStatus:  "pending",
			Enrollment:     "{}",
		}

		if result := db.Create(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create order")
		}

		if result := db.Preload("Package").First(&order, order.ID); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to load order")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Order created", order)
	}
}

func GetUserOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		query := db.Model(&models.Order{}).Where("user_id = ?", uid)
		if status := c.Query("payment_status"); status != "" {
			query = query.Where("payment_status = ?", status)
		}

		var total int64
		query.Count(&total)

		var orders []models.Order
		if result := query.Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve orders")
		}

		return utils.SuccessResponseWithPagination(c, "Orders retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetOrderByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("Package").Where("id = ? AND user_id = ?", id, uid).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Order retrieved", order)
	}
}

func SimulatePayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()

		var order models.Order
		if result := db.Where("id = ? AND user_id = ?", id, uid).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found or access denied")
		}
		if order.PaymentStatus != "pending" {
			return utils.BadRequestResponse(c, "Order is not pending", nil)
		}

		now := time.Now()
		order.PaymentStatus = "paid"
		order.TransactionID = "TRX-" + utils.GenerateRandomString(10)
		order.PaymentType = "manual"
		order.PaymentTime = &now

		if order.PackageID != nil {
			var pkg models.Package
			if result := db.Where("id = ?", order.PackageID).First(&pkg); result.Error == nil {
				dur := pkg.MembershipDurationDays
				order.MembershipDurationDays = &dur
				order.MembershipStart = &now
				end := now.AddDate(0, 0, pkg.MembershipDurationDays)
				order.MembershipEnd = &end
			}
		}

		if result := db.Save(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to process payment")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Payment successful", order)
	}
}

func PaymentStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var orders []models.Order
		if result := db.Where("user_id = ?", uid).Order("created_at DESC").Limit(10).Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve payment status")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Payment status retrieved", orders)
	}
}

func AdminIndexOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		query := db.Model(&models.Order{})
		if status := c.Query("payment_status"); status != "" {
			query = query.Where("payment_status = ?", status)
		}
		var total int64
		query.Count(&total)
		var orders []models.Order
		if result := query.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve orders")
		}
		return utils.SuccessResponseWithPagination(c, "Orders retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminShowOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Order retrieved", order)
	}
}

func AdminVerifyOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found")
		}
		order.PaymentStatus = "paid"
		now := time.Now()
		order.PaymentTime = &now
		if order.PackageID != nil {
			var pkg models.Package
			if result := db.Where("id = ?", order.PackageID).First(&pkg); result.Error == nil {
				dur := pkg.MembershipDurationDays
				order.MembershipDurationDays = &dur
				order.MembershipStart = &now
				end := now.AddDate(0, 0, pkg.MembershipDurationDays)
				order.MembershipEnd = &end
			}
		}
		if result := db.Save(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to verify order")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Order verified", order)
	}
}

func GetAllTransactions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Where("payment_status = ?", "paid").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve transactions")
		}
		return utils.SuccessResponseWithPagination(c, "Transactions retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetTransactionStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var totalRevenue float64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Select("COALESCE(SUM(total_price),0)").Scan(&totalRevenue)
		var totalTransactions int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&totalTransactions)
		var totalUsers int64
		db.Model(&models.User{}).Where("role = ?", "user").Count(&totalUsers)
		return utils.SuccessResponse(c, fiber.StatusOK, "Stats retrieved", fiber.Map{
			"total_revenue": totalRevenue, "total_transactions": totalTransactions, "total_users": totalUsers,
		})
	}
}

func ShowTransaction() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid transaction ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Transaction not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Transaction retrieved", order)
	}
}

func GetPracticeStatisticsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.PracticeSession{}).Count(&total)
		var sessions []models.PracticeSession
		if result := db.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&sessions); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve statistics")
		}
		return utils.SuccessResponseWithPagination(c, "Statistics retrieved", sessions, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func ShowPracticeStatistics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid session ID", nil)
		}
		db := database.GetDB()
		var session models.PracticeSession
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&session); result.Error != nil {
			return utils.NotFoundResponse(c, "Session not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Session retrieved", session)
	}
}

func GetEnrollKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Where("payment_status = ?", "paid").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve enroll keys")
		}
		return utils.SuccessResponseWithPagination(c, "Enroll keys retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func ShowEnrollKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Enroll key not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Enroll key retrieved", order)
	}
}

func AdminIndexReports() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve reports")
		}
		return utils.SuccessResponseWithPagination(c, "Reports retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetVideoOrdersAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.VideoOrder{}).Count(&total)
		var orders []models.VideoOrder
		if result := db.Preload("User").Preload("Video").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve video orders")
		}
		return utils.SuccessResponseWithPagination(c, "Video orders retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminGrantVideoAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.VideoOrder
		if result := db.Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Video order not found")
		}
		var video models.Video
		db.Where("id = ?", order.VideoID).First(&video)

		now := time.Now()
		order.AccessGranted = true
		order.AccessStart = &now
		duration := video.AccessDurationDays
		if duration == 0 {
			duration = 30
		}
		end := now.AddDate(0, 0, duration)
		order.AccessEnd = &end

		if result := db.Save(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to grant access")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Access granted", order)
	}
}
