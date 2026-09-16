package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserNotifications() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)

		var notifications []models.Notification
		if result := db.Where("user_id = ?", userID).Offset(offset).Limit(perPage).Order("created_at DESC").Find(&notifications); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve notifications")
		}

		return utils.SuccessResponseWithPagination(c, "Notifications retrieved", notifications, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetUnreadNotificationCount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		var count int64
		db.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&count)
		return utils.SuccessResponse(c, fiber.StatusOK, "Unread count", fiber.Map{"count": count})
	}
}

func MarkNotificationRead() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid notification ID", nil)
		}
		db := database.GetDB()
		var notif models.Notification
		if result := db.Where("id = ? AND user_id = ?", id, userID).First(&notif); result.Error != nil {
			return utils.NotFoundResponse(c, "Notification not found")
		}
		now := time.Now()
		notif.ReadAt = &now
		db.Save(&notif)
		return utils.SuccessResponse(c, fiber.StatusOK, "Notification marked as read", notif)
	}
}

func MarkAllNotificationsRead() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		now := time.Now()
		db.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Update("read_at", now)
		return utils.SuccessResponse(c, fiber.StatusOK, "All notifications marked as read", nil)
	}
}

func AdminIndexNotifications() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.Notification{}).Count(&total)

		var notifications []models.Notification
		if result := db.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&notifications); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve notifications")
		}

		return utils.SuccessResponseWithPagination(c, "Notifications retrieved", notifications, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}
