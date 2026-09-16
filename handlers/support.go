package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SubmitSupport() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.SubmitSupportRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()

		var userID *uuid.UUID
		if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
			userID = &uid
		}

		ticket := models.SupportTicket{
			SessionID: "SUPPORT-" + utils.GenerateRandomString(8),
			UserID:    userID,
			Name:      req.Name,
			Email:     req.Email,
			Question:  req.Question,
			Status:    "pending",
		}

		if result := db.Create(&ticket); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to submit support ticket")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Support ticket submitted", ticket)
	}
}

func GetSupportTickets() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.SupportTicket{}).Count(&total)

		var tickets []models.SupportTicket
		if result := db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&tickets); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve support tickets")
		}

		return utils.SuccessResponseWithPagination(c, "Support tickets retrieved", tickets, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminIndexSupport() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		query := db.Model(&models.SupportTicket{})
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}

		var total int64
		query.Count(&total)

		var tickets []models.SupportTicket
		if result := query.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&tickets); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve support tickets")
		}

		return utils.SuccessResponseWithPagination(c, "Support tickets retrieved", tickets, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminShowSupport() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ticket ID", nil)
		}
		db := database.GetDB()
		var ticket models.SupportTicket
		if result := db.Preload("User").Where("id = ?", id).First(&ticket); result.Error != nil {
			return utils.NotFoundResponse(c, "Support ticket not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Support ticket retrieved", ticket)
	}
}

func AdminAnswerSupport() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ticket ID", nil)
		}
		db := database.GetDB()
		var ticket models.SupportTicket
		if result := db.Where("id = ?", id).First(&ticket); result.Error != nil {
			return utils.NotFoundResponse(c, "Support ticket not found")
		}

		var req models.AnswerSupportRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		now := time.Now()
		ticket.Answer = req.Answer
		ticket.Status = "answered"
		ticket.AnsweredAt = &now
		db.Save(&ticket)

		return utils.SuccessResponse(c, fiber.StatusOK, "Answer submitted", ticket)
	}
}

func AdminUpdateSupportStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ticket ID", nil)
		}
		db := database.GetDB()
		var ticket models.SupportTicket
		if result := db.Where("id = ?", id).First(&ticket); result.Error != nil {
			return utils.NotFoundResponse(c, "Support ticket not found")
		}

		var req models.UpdateSupportStatusRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		ticket.Status = req.Status
		db.Save(&ticket)

		return utils.SuccessResponse(c, fiber.StatusOK, "Status updated", ticket)
	}
}

func AdminDeleteSupport() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ticket ID", nil)
		}
		db := database.GetDB()
		var ticket models.SupportTicket
		if result := db.Where("id = ?", id).First(&ticket); result.Error != nil {
			return utils.NotFoundResponse(c, "Support ticket not found")
		}
		db.Delete(&ticket)
		return utils.SuccessResponse(c, fiber.StatusOK, "Ticket deleted", nil)
	}
}
