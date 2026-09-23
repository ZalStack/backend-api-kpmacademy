package handlers

import (
	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UserDashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var totalOrders int64
		db.Model(&models.Order{}).Where("user_id = ? AND payment_status = ?", uid, "paid").Count(&totalOrders)
		var totalPractice int64
		db.Model(&models.PracticeSession{}).Where("user_id = ?", uid).Count(&totalPractice)
		var avgScore float64
		db.Model(&models.PracticeSession{}).Where("user_id = ? AND status = ?", uid, "finished").Select("COALESCE(AVG(total_score),0)").Scan(&avgScore)
		return utils.SuccessResponse(c, fiber.StatusOK, "Dashboard data", fiber.Map{
			"total_orders": totalOrders, "total_practice": totalPractice, "avg_score": avgScore,
		})
	}
}

func GetAllPackagesPublic() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		query := db.Model(&models.Package{}).Where("is_active = ?", true)
		if kelas := c.Query("kelas"); kelas != "" {
			query = query.Where("kelas = ?", kelas)
		}
		if jenjang := c.Query("jenjang"); jenjang != "" {
			query = query.Where("jenjang = ?", jenjang)
		}
		var total int64
		query.Count(&total)
		var packages []models.Package
		if result := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&packages); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve packages")
		}
		return utils.SuccessResponseWithPagination(c, "Packages retrieved", packages, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetPackageByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		db := database.GetDB()
		var pkg models.Package
		query := db.Where("id = ?", id)
		if !utils.IsAdmin(c) {
			query = query.Where("is_active = ?", true)
		}
		if result := query.First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Package retrieved", pkg)
	}
}
