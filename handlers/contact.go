package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetContactForms() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.ContactForm{}).Count(&total)

		var forms []models.ContactForm
		if result := db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&forms); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve contact forms")
		}

		return utils.SuccessResponseWithPagination(c, "Contact forms retrieved", forms, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func CreateContactForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.SubmitContactRequest
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

		form := models.ContactForm{
			UserID:    userID,
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Subject:   req.Subject,
			Message:   req.Message,
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Status:    "pending",
		}

		if result := db.Create(&form); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to submit contact form")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Contact form submitted", form)
	}
}

func GetContactFormByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Where("id = ?", id).First(&form); result.Error != nil {
			return utils.NotFoundResponse(c, "Contact form not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Contact form retrieved", form)
	}
}

func AdminIndexContactForms() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		query := db.Model(&models.ContactForm{})
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}

		var total int64
		query.Count(&total)

		var forms []models.ContactForm
		if result := query.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&forms); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve contact forms")
		}

		return utils.SuccessResponseWithPagination(c, "Contact forms retrieved", forms, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminShowContactForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Preload("User").Where("id = ?", id).First(&form); result.Error != nil {
			return utils.NotFoundResponse(c, "Contact form not found")
		}
		if form.ReadAt == nil {
			now := time.Now()
			form.ReadAt = &now
			form.Status = "read"
			db.Save(&form)
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Contact form retrieved", form)
	}
}

func AdminReplyContactForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Where("id = ?", id).First(&form); result.Error != nil {
			return utils.NotFoundResponse(c, "Contact form not found")
		}

		var req models.ReplyContactRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		now := time.Now()
		form.AdminReply = req.AdminReply
		form.Status = "replied"
		form.RepliedAt = &now
		db.Save(&form)

		return utils.SuccessResponse(c, fiber.StatusOK, "Reply submitted", form)
	}
}

func AdminUpdateContactStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Where("id = ?", id).First(&form); result.Error != nil {
			return utils.NotFoundResponse(c, "Contact form not found")
		}

		var req models.UpdateContactStatusRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		form.Status = req.Status
		db.Save(&form)
		return utils.SuccessResponse(c, fiber.StatusOK, "Status updated", form)
	}
}

func AdminDeleteContactForm() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Where("id = ?", id).First(&form); result.Error != nil {
			return utils.NotFoundResponse(c, "Contact form not found")
		}
		db.Delete(&form)
		return utils.SuccessResponse(c, fiber.StatusOK, "Contact form deleted", nil)
	}
}

func SendChatMessage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()

		var req models.SendChatRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		var userID *uuid.UUID
		if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
			userID = &uid
		}

		sessionID := req.SessionID
		if sessionID == "" {
			sessionID = "CHAT-" + utils.GenerateRandomString(8)
		}

		userMsg := models.ChatMessage{
			SessionID: sessionID,
			UserID:    userID,
			Role:      "user",
			Message:   req.Message,
			IsAI:      false,
		}
		db.Create(&userMsg)

		assistantMsg := models.ChatMessage{
			SessionID: sessionID,
			UserID:    userID,
			Role:      "assistant",
			Message:   "Thank you for your message. Our team will get back to you soon.",
			IsAI:      true,
		}
		db.Create(&assistantMsg)

		return utils.SuccessResponse(c, fiber.StatusCreated, "Message sent", fiber.Map{
			"session_id": sessionID,
			"user_message":     userMsg,
			"assistant_message": assistantMsg,
		})
	}
}

func GetChatHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Query("session_id")
		if sessionID == "" {
			return utils.BadRequestResponse(c, "session_id is required", nil)
		}
		db := database.GetDB()

		var messages []models.ChatMessage
		if result := db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&messages); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve chat history")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Chat history retrieved", messages)
	}
}
