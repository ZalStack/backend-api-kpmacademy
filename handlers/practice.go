package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func StartPractice() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()

		var req models.StartPracticeRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		var pkg models.Package
		if result := db.Where("id = ?", req.PackageID).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}

		var order models.Order
		if result := db.Where("user_id = ? AND package_id = ? AND payment_status = ?", userID, pkg.ID, "paid").First(&order); result.Error != nil {
			return utils.ForbiddenResponse(c, "You must purchase this package first")
		}

		now := time.Now()
		session := models.PracticeSession{
			UserID:    userID,
			PackageID: pkg.ID,
			OrderID:   order.ID,
			CardID:    req.CardID,
			StartedAt: &now,
			Status:    "in_progress",
		}

		if result := db.Create(&session); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to start practice session")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Practice started", session)
	}
}

func SubmitPractice() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()

		var req models.SubmitPracticeRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		var session models.PracticeSession
		if result := db.Where("id = ? AND user_id = ?", req.SessionID, userID).First(&session); result.Error != nil {
			return utils.NotFoundResponse(c, "Practice session not found")
		}
		if session.Status == "finished" {
			return utils.BadRequestResponse(c, "Session already finished", nil)
		}

		now := time.Now()
		session.Answers = req.Answers
		session.DurationSeconds = req.Duration
		session.FinishedAt = &now
		session.Status = "finished"

		db.Save(&session)
		return utils.SuccessResponse(c, fiber.StatusOK, "Practice submitted", session)
	}
}

func ShowPracticeSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid session ID", nil)
		}
		db := database.GetDB()
		var session models.PracticeSession
		if result := db.Preload("Package").Where("id = ? AND user_id = ?", id, userID).First(&session); result.Error != nil {
			return utils.NotFoundResponse(c, "Session not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Session retrieved", session)
	}
}

func GetPracticeHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.PracticeSession{}).Where("user_id = ?", userID).Count(&total)

		var sessions []models.PracticeSession
		if result := db.Preload("Package").Where("user_id = ?", userID).Offset(offset).Limit(perPage).Order("created_at DESC").Find(&sessions); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve practice history")
		}

		return utils.SuccessResponseWithPagination(c, "Practice history retrieved", sessions, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetPracticeStatisticsUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()

		var totalSessions int64
		db.Model(&models.PracticeSession{}).Where("user_id = ?", userID).Count(&totalSessions)

		var avgScore float64
		db.Model(&models.PracticeSession{}).Where("user_id = ? AND status = ?", userID, "finished").Select("COALESCE(AVG(total_score),0)").Scan(&avgScore)

		var totalCorrect int64
		db.Model(&models.PracticeSession{}).Where("user_id = ? AND status = ?", userID, "finished").Select("COALESCE(SUM(correct_answer),0)").Scan(&totalCorrect)

		var totalWrong int64
		db.Model(&models.PracticeSession{}).Where("user_id = ? AND status = ?", userID, "finished").Select("COALESCE(SUM(wrong_answer),0)").Scan(&totalWrong)

		return utils.SuccessResponse(c, fiber.StatusOK, "Statistics retrieved", fiber.Map{
			"total_sessions": totalSessions,
			"avg_score":      avgScore,
			"total_correct":  totalCorrect,
			"total_wrong":    totalWrong,
		})
	}
}
