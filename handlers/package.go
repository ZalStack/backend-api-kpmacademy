package handlers

import (
	"encoding/json"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var jsonUnmarshal = json.Unmarshal
var jsonMarshal = json.Marshal

func AdminDashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var totalUsers int64
		db.Model(&models.User{}).Where("role = ?", "user").Count(&totalUsers)
		var totalPackages int64
		db.Model(&models.Package{}).Count(&totalPackages)
		var totalOrders int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&totalOrders)
		var totalRevenue float64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Select("COALESCE(SUM(total_price),0)").Scan(&totalRevenue)
		return utils.SuccessResponse(c, fiber.StatusOK, "Dashboard data", fiber.Map{
			"total_users": totalUsers, "total_packages": totalPackages,
			"total_orders": totalOrders, "total_revenue": totalRevenue,
		})
	}
}

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

func GetAllPackagesAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Package{}).Count(&total)
		var packages []models.Package
		if result := db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&packages); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve packages")
		}
		return utils.SuccessResponseWithPagination(c, "Packages retrieved", packages, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
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

func CreatePackage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.CreatePackageRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}
		db := database.GetDB()
		cards := req.Cards
		if cards == "" {
			cards = "[]"
		}
		questions := req.Questions
		if questions == "" {
			questions = "[]"
		}
		pkg := models.Package{
			Title: req.Title, Description: req.Description, Thumbnail: req.Thumbnail,
			Kelas: req.Kelas, Jenjang: req.Jenjang, Price: req.Price,
			IsPayWhatYouWant: req.IsPayWhatYouWant, MinPayAmount: req.MinPayAmount,
			MembershipDurationDays: req.MembershipDurationDays, Cards: cards,
			Questions: questions, Reviews: "[]", HideExplanation: req.HideExplanation, IsActive: true,
		}
		if pkg.MembershipDurationDays == 0 {
			pkg.MembershipDurationDays = 30
		}
		if req.DiscountPrice > 0 {
			pkg.DiscountPrice = &req.DiscountPrice
			pkg.IsDiscountActive = req.IsDiscountActive
		}
		if req.TimeLimitMinutes > 0 {
			pkg.TimeLimitMinutes = &req.TimeLimitMinutes
		}
		if result := db.Create(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create package")
		}
		return utils.SuccessResponse(c, fiber.StatusCreated, "Package created", pkg)
	}
}

func UpdatePackage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		db := database.GetDB()
		var pkg models.Package
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}
		var req models.UpdatePackageRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if req.Title != "" {
			pkg.Title = req.Title
		}
		if req.Description != "" {
			pkg.Description = req.Description
		}
		if req.Thumbnail != "" {
			pkg.Thumbnail = req.Thumbnail
		}
		if req.Kelas != "" {
			pkg.Kelas = req.Kelas
		}
		if req.Jenjang != "" {
			pkg.Jenjang = req.Jenjang
		}
		if req.Price != nil {
			pkg.Price = *req.Price
		}
		if req.DiscountPrice != nil {
			pkg.DiscountPrice = req.DiscountPrice
		}
		if req.IsDiscountActive != nil {
			pkg.IsDiscountActive = *req.IsDiscountActive
		}
		if req.IsPayWhatYouWant != nil {
			pkg.IsPayWhatYouWant = *req.IsPayWhatYouWant
		}
		if req.MinPayAmount != nil {
			pkg.MinPayAmount = *req.MinPayAmount
		}
		if req.MembershipDurationDays != nil {
			pkg.MembershipDurationDays = *req.MembershipDurationDays
		}
		if req.Cards != "" {
			pkg.Cards = req.Cards
		}
		if req.Questions != "" {
			pkg.Questions = req.Questions
		}
		if req.HideExplanation != nil {
			pkg.HideExplanation = *req.HideExplanation
		}
		if req.TimeLimitMinutes != nil {
			pkg.TimeLimitMinutes = req.TimeLimitMinutes
		}
		if req.IsActive != nil {
			pkg.IsActive = *req.IsActive
		}
		if result := db.Save(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to update package")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Package updated", pkg)
	}
}

func DeletePackage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		db := database.GetDB()
		var pkg models.Package
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}
		if result := db.Delete(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to delete package")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Package deleted", nil)
	}
}

func AddCardToPackage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		db := database.GetDB()
		var pkg models.Package
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}
		var req models.AddCardRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		pkg.Cards = req.Cards
		if result := db.Save(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to update cards")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Cards updated", pkg)
	}
}

func RemoveCardFromPackage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		cardID := c.Params("cardId")
		if cardID == "" {
			return utils.BadRequestResponse(c, "Invalid card ID", nil)
		}

		db := database.GetDB()
		var pkg models.Package
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}

		if pkg.Cards == "" || pkg.Cards == "[]" || pkg.Cards == "null" {
			return utils.SuccessResponse(c, fiber.StatusOK, "No cards to remove", pkg)
		}

		var cards []map[string]interface{}
		if err := jsonUnmarshal([]byte(pkg.Cards), &cards); err != nil {
			return utils.InternalErrorResponse(c, "Failed to parse cards")
		}

		var filtered []map[string]interface{}
		found := false
		for _, card := range cards {
			cid, _ := card["id"].(string)
			if cid == cardID {
				found = true
				continue
			}
			filtered = append(filtered, card)
		}

		if !found {
			return utils.NotFoundResponse(c, "Card not found")
		}

		filteredBytes, err := jsonMarshal(filtered)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to serialize cards")
		}
		pkg.Cards = string(filteredBytes)

		if result := db.Save(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to remove card")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Card removed", pkg)
	}
}

func ImportQuestionsPDF() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid package ID", nil)
		}
		db := database.GetDB()
		var pkg models.Package
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
			return utils.NotFoundResponse(c, "Package not found")
		}
		var req models.ImportPDFRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		pkg.Questions = req.Questions
		if result := db.Save(&pkg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to import questions")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Questions imported", pkg)
	}
}
