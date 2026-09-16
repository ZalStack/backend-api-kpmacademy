package handlers

import (
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
		userID := c.Locals("user_id").(uuid.UUID)
		db := database.GetDB()
		var totalOrders int64
		db.Model(&models.Order{}).Where("user_id = ? AND payment_status = ?", userID, "paid").Count(&totalOrders)
		var totalPractice int64
		db.Model(&models.PracticeSession{}).Where("user_id = ?", userID).Count(&totalPractice)
		var avgScore float64
		db.Model(&models.PracticeSession{}).Where("user_id = ? AND status = ?", userID, "finished").Select("COALESCE(AVG(total_score),0)").Scan(&avgScore)
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
		if result := db.Where("id = ?", id).First(&pkg); result.Error != nil {
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
		pkg := models.Package{
			Title: req.Title, Description: req.Description, Thumbnail: req.Thumbnail,
			Kelas: req.Kelas, Jenjang: req.Jenjang, Price: req.Price,
			IsPayWhatYouWant: req.IsPayWhatYouWant, MinPayAmount: req.MinPayAmount,
			MembershipDurationDays: req.MembershipDurationDays, Cards: req.Cards,
			Questions: req.Questions, HideExplanation: req.HideExplanation, IsActive: true,
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
		db.Save(&pkg)
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
		db.Delete(&pkg)
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
		db.Save(&pkg)
		return utils.SuccessResponse(c, fiber.StatusOK, "Cards updated", pkg)
	}
}

func RemoveCardFromPackage() fiber.Handler {
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
		pkg.Cards = "[]"
		db.Save(&pkg)
		return utils.SuccessResponse(c, fiber.StatusOK, "Cards cleared", pkg)
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
		db.Save(&pkg)
		return utils.SuccessResponse(c, fiber.StatusOK, "Questions imported", pkg)
	}
}

func AdminIndexOrders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		query := db.Model(&models.Order{})
		if status := c.Query("payment_status"); status != "" {
			query = query.Where("payment_status = ?", status)
		}
		var total int64
		query.Count(&total)
		var orders []models.Order
		if result := query.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve orders")
		}
		return utils.SuccessResponseWithPagination(c, "Orders retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminShowOrder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Order not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Order retrieved", order)
	}
}

func AdminVerifyOrder() fiber.Handler {
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
		order.PaymentStatus = "paid"
		now := time.Now()
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
		return utils.SuccessResponse(c, fiber.StatusOK, "Order verified", order)
	}
}

func GetAllTransactions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Where("payment_status = ?", "paid").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve transactions")
		}
		return utils.SuccessResponseWithPagination(c, "Transactions retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetTransactionStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var totalRevenue float64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Select("COALESCE(SUM(total_price),0)").Scan(&totalRevenue)
		var totalTransactions int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&totalTransactions)
		var totalUsers int64
		db.Model(&models.User{}).Where("role = ?", "user").Count(&totalUsers)
		return utils.SuccessResponse(c, fiber.StatusOK, "Stats retrieved", fiber.Map{
			"total_revenue": totalRevenue, "total_transactions": totalTransactions, "total_users": totalUsers,
		})
	}
}

func ShowTransaction() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid transaction ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Transaction not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Transaction retrieved", order)
	}
}

func GetPracticeStatisticsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.PracticeSession{}).Count(&total)
		var sessions []models.PracticeSession
		if result := db.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&sessions); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve statistics")
		}
		return utils.SuccessResponseWithPagination(c, "Statistics retrieved", sessions, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func ShowPracticeStatistics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid session ID", nil)
		}
		db := database.GetDB()
		var session models.PracticeSession
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&session); result.Error != nil {
			return utils.NotFoundResponse(c, "Session not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Session retrieved", session)
	}
}

func GetEnrollKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Where("payment_status = ?", "paid").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve enroll keys")
		}
		return utils.SuccessResponseWithPagination(c, "Enroll keys retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func ShowEnrollKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid ID", nil)
		}
		db := database.GetDB()
		var order models.Order
		if result := db.Preload("User").Preload("Package").Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Enroll key not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Enroll key retrieved", order)
	}
}

func AdminIndexReports() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.Order{}).Count(&total)
		var orders []models.Order
		if result := db.Preload("User").Preload("Package").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve reports")
		}
		return utils.SuccessResponseWithPagination(c, "Reports retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetVideoOrdersAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)
		var total int64
		db.Model(&models.VideoOrder{}).Count(&total)
		var orders []models.VideoOrder
		if result := db.Preload("User").Preload("Video").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&orders); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve video orders")
		}
		return utils.SuccessResponseWithPagination(c, "Video orders retrieved", orders, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func AdminGrantVideoAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid order ID", nil)
		}
		db := database.GetDB()
		var order models.VideoOrder
		if result := db.Where("id = ?", id).First(&order); result.Error != nil {
			return utils.NotFoundResponse(c, "Video order not found")
		}
		now := time.Now()
		order.AccessGranted = true
		order.AccessStart = &now
		end := now.AddDate(0, 0, 30)
		order.AccessEnd = &end
		db.Save(&order)
		return utils.SuccessResponse(c, fiber.StatusOK, "Access granted", order)
	}
}
