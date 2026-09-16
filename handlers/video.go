package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllVideos() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		query := db.Model(&models.Video{})
		if utils.IsAdmin(c) {
			// Admin sees all
		} else {
			query = query.Where("is_active = ?", true)
		}
		if packageID := c.Query("package_id"); packageID != "" {
			query = query.Where("package_id = ?", packageID)
		}

		var total int64
		query.Count(&total)

		var videos []models.Video
		if result := query.Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&videos); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve videos")
		}

		return utils.SuccessResponseWithPagination(c, "Videos retrieved", videos, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetVideoByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid video ID", nil)
		}
		db := database.GetDB()
		var video models.Video
		if result := db.Preload("Package").Where("id = ?", id).First(&video); result.Error != nil {
			return utils.NotFoundResponse(c, "Video not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Video retrieved", video)
	}
}

func CreateVideo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.CreateVideoRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()
		video := models.Video{
			Title:              req.Title,
			Description:        req.Description,
			VideoFile:          req.VideoFile,
			VideoURL:           req.VideoURL,
			Thumbnail:          req.Thumbnail,
			Price:              req.Price,
			DiscountType:       req.DiscountType,
			AccessDurationDays: req.AccessDurationDays,
			IsPayWhatYouWant:   req.IsPayWhatYouWant,
			MinPayAmount:       req.MinPayAmount,
			IsActive:           true,
		}
		if video.AccessDurationDays == 0 {
			video.AccessDurationDays = 30
		}
		if req.DiscountValue > 0 {
			video.DiscountValue = &req.DiscountValue
		}
		if req.PackageID != "" {
			pid, err := uuid.Parse(req.PackageID)
			if err == nil {
				video.PackageID = &pid
			}
		}

		if result := db.Create(&video); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create video")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Video created", video)
	}
}

func UpdateVideo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid video ID", nil)
		}
		db := database.GetDB()
		var video models.Video
		if result := db.Where("id = ?", id).First(&video); result.Error != nil {
			return utils.NotFoundResponse(c, "Video not found")
		}

		var req models.UpdateVideoRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		if req.Title != "" {
			video.Title = req.Title
		}
		if req.Description != "" {
			video.Description = req.Description
		}
		if req.VideoFile != "" {
			video.VideoFile = req.VideoFile
		}
		if req.VideoURL != "" {
			video.VideoURL = req.VideoURL
		}
		if req.Thumbnail != "" {
			video.Thumbnail = req.Thumbnail
		}
		if req.Price != nil {
			video.Price = *req.Price
		}
		if req.DiscountType != "" {
			video.DiscountType = req.DiscountType
		}
		if req.DiscountValue != nil {
			video.DiscountValue = req.DiscountValue
		}
		if req.AccessDurationDays != nil {
			video.AccessDurationDays = *req.AccessDurationDays
		}
		if req.IsActive != nil {
			video.IsActive = *req.IsActive
		}
		if req.IsPayWhatYouWant != nil {
			video.IsPayWhatYouWant = *req.IsPayWhatYouWant
		}
		if req.MinPayAmount != nil {
			video.MinPayAmount = *req.MinPayAmount
		}

		if result := db.Save(&video); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to update video")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Video updated", video)
	}
}

func DeleteVideo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid video ID", nil)
		}
		db := database.GetDB()
		var video models.Video
		if result := db.Where("id = ?", id).First(&video); result.Error != nil {
			return utils.NotFoundResponse(c, "Video not found")
		}
		if result := db.Delete(&video); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to delete video")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Video deleted", nil)
	}
}

func ToggleVideoActive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid video ID", nil)
		}
		db := database.GetDB()
		var video models.Video
		if result := db.Where("id = ?", id).First(&video); result.Error != nil {
			return utils.NotFoundResponse(c, "Video not found")
		}
		video.IsActive = !video.IsActive
		if result := db.Save(&video); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to toggle video status")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Video status toggled", video)
	}
}

func CreateVideoOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		videoID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid video ID", nil)
		}
		db := database.GetDB()

		var video models.Video
		if result := db.Where("id = ? AND is_active = ?", videoID, true).First(&video); result.Error != nil {
			return utils.NotFoundResponse(c, "Video not found")
		}

		totalPrice := video.Price
		if video.IsPayWhatYouWant {
			var req models.CreateVideoOrderRequest
			if err := c.BodyParser(&req); err == nil && req.TotalPrice > 0 {
				if req.TotalPrice < video.MinPayAmount {
					return utils.BadRequestResponse(c, "Minimum payment is required", fiber.Map{
						"min_amount": video.MinPayAmount,
					})
				}
				totalPrice = req.TotalPrice
			}
		}

		orderNumber := "VID-ORD-" + time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(6)

		order := models.VideoOrder{
			UserID:        uid,
			VideoID:       videoID,
			OrderNumber:   orderNumber,
			TotalPrice:    totalPrice,
			PaymentStatus: "pending",
		}

		if result := db.Create(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create video order")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Video order created", order)
	}
}

func SimulateVideoPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		id, err := uuid.Parse(c.Params("orderId"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()

		var order models.VideoOrder
		if result := db.Where("id = ? AND user_id = ?", id, uid).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Video order not found or access denied")
		}
		if order.PaymentStatus != "pending" {
			return utils.BadRequestResponse(c, "Order is not pending", nil)
		}

		now := time.Now()
		order.PaymentStatus = "paid"
		order.TransactionID = "VID-TRX-" + utils.GenerateRandomString(10)
		order.PaymentType = "manual"
		order.PaymentTime = &now
		order.AccessGranted = true
		order.AccessStart = &now

		var video models.Video
		if result := db.Where("id = ?", order.VideoID).First(&video); result.Error == nil {
			duration := video.AccessDurationDays
			if duration == 0 {
				duration = 30
			}
			end := now.AddDate(0, 0, duration)
			order.AccessEnd = &end
		} else {
			end := now.AddDate(0, 0, 30)
			order.AccessEnd = &end
		}

		if result := db.Save(&order); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to process payment")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Video payment successful", order)
	}
}
