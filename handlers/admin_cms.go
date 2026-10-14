package handlers

import (
	"net/http"
	"strings"
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ==========================================
// BANNER MANAGEMENT (ADMIN)
// ==========================================

type CreateBannerRequest struct {
	Title     string `json:"title" validate:"required"`
	ImageURL  string `json:"image_url" validate:"required"`
	Route     string `json:"route"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
	Tag       string `json:"tag"`
	SubTag    string `json:"sub_tag"`
	Subtitle  string `json:"subtitle"`
}

func CreateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateBannerRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload banner tidak valid", err.Error())
		}
		if req.Title == "" || req.ImageURL == "" {
			return utils.BadRequestResponse(c, "Judul dan URL gambar banner wajib diisi", nil)
		}

		db := database.GetDB()
		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
		route := req.Route
		if route == "" {
			route = "/package_list"
		}

		banner := models.Banner{
			ID:        uuid.New(),
			Title:     req.Title,
			ImageURL:  req.ImageURL,
			Route:     route,
			SortOrder: req.SortOrder,
			IsActive:  isActive,
			Tag:       req.Tag,
			SubTag:    req.SubTag,
			Subtitle:  req.Subtitle,
		}

		if err := db.Create(&banner).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal menyimpan banner: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusCreated, "Banner berhasil dibuat", banner)
	}
}

func UpdateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID banner tidak valid", nil)
		}

		db := database.GetDB()
		var banner models.Banner
		if err := db.Where("id = ?", id).First(&banner).Error; err != nil {
			return utils.NotFoundResponse(c, "Banner tidak ditemukan")
		}

		var req CreateBannerRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload tidak valid", err.Error())
		}

		if req.Title != "" {
			banner.Title = req.Title
		}
		if req.ImageURL != "" {
			banner.ImageURL = req.ImageURL
		}
		if req.Route != "" {
			banner.Route = req.Route
		}
		if req.SortOrder != 0 {
			banner.SortOrder = req.SortOrder
		}
		if req.IsActive != nil {
			banner.IsActive = *req.IsActive
		}

		if err := db.Save(&banner).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal memperbarui banner: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Banner berhasil diperbarui", banner)
	}
}

func DeleteBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID banner tidak valid", nil)
		}

		db := database.GetDB()
		if err := db.Where("id = ?", id).Delete(&models.Banner{}).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal menghapus banner: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Banner berhasil dihapus", nil)
	}
}

// ==========================================
// LIVE CLASS MANAGEMENT (ADMIN)
// ==========================================

type CreateLiveClassRequest struct {
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description"`
	TutorName       string `json:"tutor_name" validate:"required"`
	TutorRole       string `json:"tutor_role"`
	TutorAvatar     string `json:"tutor_avatar"`
	Subject         string `json:"subject"`
	Grade           string `json:"grade"`
	ScheduledAt     string `json:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes"`
	MeetingURL      string `json:"meeting_url"`
	ZoomURL         string `json:"zoom_url"`
	IsLive          *bool  `json:"is_live"`
	IsActive        *bool  `json:"is_active"`
}

func CreateLiveClass() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateLiveClassRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload live class tidak valid", err.Error())
		}
		if req.Title == "" {
			return utils.BadRequestResponse(c, "Topik live class wajib diisi", nil)
		}

		db := database.GetDB()
		zoomURL := req.ZoomURL
		if zoomURL == "" {
			zoomURL = req.MeetingURL
		}
		if zoomURL == "" {
			zoomURL = "https://zoom.us"
		}

		tutorName := req.TutorName
		if tutorName == "" {
			tutorName = "Tutor KPM"
		}

		sched := time.Now().Add(24 * time.Hour)
		if req.ScheduledAt != "" {
			if parsed, err := time.Parse(time.RFC3339, req.ScheduledAt); err == nil {
				sched = parsed
			}
		}

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		live := models.LiveClass{
			ID:              uuid.New(),
			Title:           req.Title,
			Description:     req.Description,
			TutorName:       tutorName,
			TutorRole:       req.TutorRole,
			TutorAvatar:     req.TutorAvatar,
			Subject:         req.Subject,
			Grade:           req.Grade,
			ScheduledAt:     sched,
			DurationMinutes: req.DurationMinutes,
			ZoomURL:         zoomURL,
			IsActive:        isActive,
		}
		if live.DurationMinutes <= 0 {
			live.DurationMinutes = 60
		}

		if err := db.Create(&live).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal membuat live class: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusCreated, "Jadwal live class berhasil dibuat", live)
	}
}

func UpdateLiveClass() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID live class tidak valid", nil)
		}

		db := database.GetDB()
		var live models.LiveClass
		if err := db.Where("id = ?", id).First(&live).Error; err != nil {
			return utils.NotFoundResponse(c, "Live class tidak ditemukan")
		}

		var req CreateLiveClassRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload tidak valid", err.Error())
		}

		if req.Title != "" {
			live.Title = req.Title
		}
		if req.TutorName != "" {
			live.TutorName = req.TutorName
		}
		if req.ZoomURL != "" {
			live.ZoomURL = req.ZoomURL
		} else if req.MeetingURL != "" {
			live.ZoomURL = req.MeetingURL
		}
		if req.Subject != "" {
			live.Subject = req.Subject
		}
		if req.Grade != "" {
			live.Grade = req.Grade
		}
		if req.IsActive != nil {
			live.IsActive = *req.IsActive
		}

		if err := db.Save(&live).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal memperbarui live class: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Live class berhasil diperbarui", live)
	}
}

func DeleteLiveClass() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID live class tidak valid", nil)
		}

		db := database.GetDB()
		if err := db.Where("id = ?", id).Delete(&models.LiveClass{}).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal menghapus live class: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Live class berhasil dihapus", nil)
	}
}

// ==========================================
// MODULE MANAGEMENT (ADMIN)
// ==========================================

type CreateModuleRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Subject     string  `json:"subject"`
	Grade       string  `json:"grade"`
	Chapter     int     `json:"chapter"`
	TotalPages  int     `json:"total_pages"`
	FileSizeMb  float64 `json:"file_size_mb"`
	PDFURL      string  `json:"pdf_url"`
	FileURL     string  `json:"file_url"`
	CoverURL    string  `json:"cover_url"`
	IsActive    *bool   `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

func CreateModule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateModuleRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload modul tidak valid", err.Error())
		}
		if req.Title == "" {
			return utils.BadRequestResponse(c, "Judul modul wajib diisi", nil)
		}

		db := database.GetDB()
		pdfURL := req.PDFURL
		if pdfURL == "" {
			pdfURL = req.FileURL
		}
		if pdfURL == "" {
			pdfURL = "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
		}

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		coverImg := req.CoverURL
		mod := models.Module{
			ID:          uuid.New(),
			Title:       req.Title,
			Description: req.Description,
			Subject:     req.Subject,
			Grade:       req.Grade,
			Chapter:     req.Chapter,
			TotalPages:  req.TotalPages,
			PDFURL:      pdfURL,
			CoverImage:  coverImg,
			IsActive:    isActive,
			SortOrder:   req.SortOrder,
		}
		if mod.TotalPages <= 0 {
			mod.TotalPages = 20
		}

		if err := db.Create(&mod).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal membuat modul: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusCreated, "Modul PDF berhasil dibuat", mod)
	}
}

func UpdateModule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID modul tidak valid", nil)
		}

		db := database.GetDB()
		var mod models.Module
		if err := db.Where("id = ?", id).First(&mod).Error; err != nil {
			return utils.NotFoundResponse(c, "Modul tidak ditemukan")
		}

		var req CreateModuleRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload tidak valid", err.Error())
		}

		if req.Title != "" {
			mod.Title = req.Title
		}
		if req.Subject != "" {
			mod.Subject = req.Subject
		}
		if req.Grade != "" {
			mod.Grade = req.Grade
		}
		if req.CoverURL != "" {
			mod.CoverImage = req.CoverURL
		}
		if req.PDFURL != "" {
			mod.PDFURL = req.PDFURL
		} else if req.FileURL != "" {
			mod.PDFURL = req.FileURL
		}
		if req.IsActive != nil {
			mod.IsActive = *req.IsActive
		}

		if err := db.Save(&mod).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal memperbarui modul: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Modul berhasil diperbarui", mod)
	}
}

func DeleteModule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID modul tidak valid", nil)
		}

		db := database.GetDB()
		if err := db.Where("id = ?", id).Delete(&models.Module{}).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal menghapus modul: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Modul berhasil dihapus", nil)
	}
}

// ==========================================
// VOUCHER MANAGEMENT (ADMIN)
// ==========================================

type CreateVoucherRequest struct {
	Code            string  `json:"code" validate:"required"`
	Title           string  `json:"title" validate:"required"`
	Description     string  `json:"description"`
	DiscountPercent int     `json:"discount_percent"`
	DiscountAmount  float64 `json:"discount_amount"`
	MinPurchase     float64 `json:"min_purchase"`
	MaxDiscount     float64 `json:"max_discount"`
	Quota           int     `json:"quota"`
	UsageLimit      int     `json:"usage_limit"`
}

func GetAllVouchersAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var vouchers []models.Voucher
		db.Order("created_at desc").Find(&vouchers)
		return utils.SuccessResponse(c, http.StatusOK, "Daftar voucher berhasil dimuat", vouchers)
	}
}

func CreateVoucher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateVoucherRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload voucher tidak valid", err.Error())
		}
		code := strings.ToUpper(strings.TrimSpace(req.Code))
		if code == "" || req.Title == "" {
			return utils.BadRequestResponse(c, "Kode dan judul voucher wajib diisi", nil)
		}

		quota := req.Quota
		if quota <= 0 {
			quota = req.UsageLimit
		}
		if quota <= 0 {
			quota = 100
		}

		db := database.GetDB()
		voucher := models.Voucher{
			ID:              uuid.New(),
			Code:            code,
			Title:           req.Title,
			Description:     req.Description,
			DiscountPercent: req.DiscountPercent,
			DiscountAmount:  req.DiscountAmount,
			MinPurchase:     req.MinPurchase,
			MaxDiscount:     req.MaxDiscount,
			ValidUntil:      time.Now().AddDate(0, 3, 0),
			IsActive:        true,
			UsageLimit:      quota,
			UsageCount:      0,
		}

		if err := db.Create(&voucher).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal membuat voucher: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusCreated, "Voucher berhasil dibuat", voucher)
	}
}

func UpdateVoucher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID voucher tidak valid", nil)
		}

		db := database.GetDB()
		var voucher models.Voucher
		if err := db.Where("id = ?", id).First(&voucher).Error; err != nil {
			return utils.NotFoundResponse(c, "Voucher tidak ditemukan")
		}

		var req CreateVoucherRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload tidak valid", err.Error())
		}

		if req.Code != "" {
			voucher.Code = strings.ToUpper(strings.TrimSpace(req.Code))
		}
		if req.Title != "" {
			voucher.Title = req.Title
		}
		if req.DiscountPercent > 0 {
			voucher.DiscountPercent = req.DiscountPercent
		}
		if req.DiscountAmount > 0 {
			voucher.DiscountAmount = req.DiscountAmount
		}
		if req.Quota > 0 {
			voucher.UsageLimit = req.Quota
		}

		if err := db.Save(&voucher).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal memperbarui voucher: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Voucher berhasil diperbarui", voucher)
	}
}

func DeleteVoucher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID voucher tidak valid", nil)
		}

		db := database.GetDB()
		if err := db.Where("id = ?", id).Delete(&models.Voucher{}).Error; err != nil {
			return utils.InternalErrorResponse(c, "Gagal menghapus voucher: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusOK, "Voucher berhasil dihapus", nil)
	}
}
