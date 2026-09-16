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
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()

		var req models.CreateOrderRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		var pkg models.Package
		if result := db.Where("id = ? AND is_active = ?", req.PackageID, true).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found or inactive")
		}

		totalPrice := pkg.Price
		if req.TotalPrice > 0 {
			totalPrice = req.TotalPrice
		}
		if pkg.IsDiscountActive && pkg.DiscountPrice != nil {
			totalPrice = *pkg.DiscountPrice
		}

		orderNumber := "ORD-" + time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(6)
		pkgID := pkg.ID

		order := models.Order{
			UserID:        userID,
			PackageID:     &pkgID,
			Type:          "package",
			OrderNumber:   orderNumber,
			TotalPrice:    totalPrice,
			IsCustomAmount: req.IsCustomAmount,
			PaymentStatus: "pending",
		}

		if result := db.Create(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create order")
		}

		db.Preload("Package").First(&order, order.ID)
		return utils.SuccessResponse(c, fiber.StatusCreated, "Order created", order)
	}
}

func GetUserOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		query := db.Model(&models.Order{}).Where("user_id = ?", userID)
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
		userID := c.Locals("user_id").(uuid.UUID)
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("Package").Where("id = ? AND user_id = ?", id, userID).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Order retrieved", order)
	}
}

func SimulatePayment() fiber.Handler {
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

		db.Save(&order)
		return utils.SuccessResponse(c, fiber.StatusOK, "Payment successful", order)
	}
}

func PaymentStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		var orders []models.Order
		if result := db.Where("user_id = ?", userID).Order("created_at DESC").Limit(10).Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve payment status")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Payment status retrieved", orders)
	}
}
