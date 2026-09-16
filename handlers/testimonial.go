package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPublicTestimonials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.Testimonial{}).Where("is_approved = ? AND is_active = ?", true, true).Count(&total)

		var testimonials []models.Testimonial
		if result := db.Preload("User").Where("is_approved = ? AND is_active = ?", true, true).
			Offset(offset).Limit(perPage).Order("created_at DESC").Find(&testimonials); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve testimonials")
		}

		var avgRating float64
		db.Model(&models.Testimonial{}).Where("is_approved = ? AND is_active = ?", true, true).
			Select("COALESCE(AVG(rating),0)").Scan(&avgRating)

		return utils.SuccessResponseWithPagination(c, "Testimonials retrieved", fiber.Map{
			"testimonials": testimonials,
			"avg_rating":   avgRating,
			"total":        total,
		}, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func CreateTestimonial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()

		var req models.CreateTestimonialRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		rating := req.Rating
		if rating < 1 || rating > 5 {
			rating = 5
		}

		testimonial := models.Testimonial{
			UserID:  uid,
			Content: req.Content,
			Rating:  rating,
		}

		if result := db.Create(&testimonial); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create testimonial")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Testimonial submitted", testimonial)
	}
}

func GetUserTestimonial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var testimonial models.Testimonial
		if result := db.Where("user_id = ?", uid).Order("created_at DESC").First(&testimonial); result.Error != nil {
			return utils.NotFoundResponse(c, "No testimonial found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Testimonial retrieved", testimonial)
	}
}

func AdminIndexTestimonials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.Testimonial{}).Count(&total)

		var testimonials []models.Testimonial
		if result := db.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&testimonials); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve testimonials")
		}

		return utils.SuccessResponseWithPagination(c, "Testimonials retrieved", testimonials, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminApproveTestimonial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid testimonial ID", nil)
		}
		db := database.GetDB()
		var testimonial models.Testimonial
		if result := db.Where("id = ?", id).First(&testimonial); result.Error != nil {
			return utils.NotFoundResponse(c, "Testimonial not found")
		}

		now := time.Now()
		testimonial.IsApproved = true
		testimonial.ApprovedAt = &now
		if result := db.Save(&testimonial); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to approve testimonial")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Testimonial approved", testimonial)
	}
}

func AdminToggleTestimonialActive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid testimonial ID", nil)
		}
		db := database.GetDB()
		var testimonial models.Testimonial
		if result := db.Where("id = ?", id).First(&testimonial); result.Error != nil {
			return utils.NotFoundResponse(c, "Testimonial not found")
		}
		testimonial.IsActive = !testimonial.IsActive
		if result := db.Save(&testimonial); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to toggle testimonial status")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Testimonial status toggled", testimonial)
	}
}

func AdminDeleteTestimonial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid testimonial ID", nil)
		}
		db := database.GetDB()
		var testimonial models.Testimonial
		if result := db.Where("id = ?", id).First(&testimonial); result.Error != nil {
			return utils.NotFoundResponse(c, "Testimonial not found")
		}
		if result := db.Delete(&testimonial); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to delete testimonial")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Testimonial deleted", nil)
	}
}

func AdminBulkDeleteTestimonials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.BulkDeleteRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if len(req.IDs) == 0 {
			return utils.BadRequestResponse(c, "IDs are required", nil)
		}

		var uuids []uuid.UUID
		for _, idStr := range req.IDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return utils.BadRequestResponse(c, "Invalid ID: "+idStr, nil)
			}
			uuids = append(uuids, id)
		}

		db := database.GetDB()
		if result := db.Where("id IN ?", uuids).Delete(&models.Testimonial{}); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to delete testimonials")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Testimonials deleted", nil)
	}
}
