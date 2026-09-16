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
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.Notification{}).Where("user_id = ?", uid).Count(&total)

		var notifications []models.Notification
		if result := db.Where("user_id = ?", uid).Offset(offset).Limit(perPage).Order("created_at DESC").Find(&notifications); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve notifications")
		}

		return utils.SuccessResponseWithPagination(c, "Notifications retrieved", notifications, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetUnreadNotificationCount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var count int64
		db.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", uid).Count(&count)
		return utils.SuccessResponse(c, fiber.StatusOK, "Unread count", fiber.Map{"count": count})
	}
}

func MarkNotificationRead() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid notification ID", nil)
		}
		db := database.GetDB()
		var notif models.Notification
		if result := db.Where("id = ? AND user_id = ?", id, uid).First(&notif); result.Error != nil {
			return utils.NotFoundResponse(c, "Notification not found")
		}
		now := time.Now()
		notif.ReadAt = &now
		if result := db.Save(&notif); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to mark notification as read")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Notification marked as read", notif)
	}
}

func MarkAllNotificationsRead() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		now := time.Now()
		if result := db.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", uid).Update("read_at", now); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to mark all notifications as read")
		}
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
