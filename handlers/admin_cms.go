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

// ResetOfficialCatalog cleans up dummy catalog data and populates standardized official KPM content
func ResetOfficialCatalog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()

		// 1. Clean existing CMS tables
		db.Exec("DELETE FROM banners")
		db.Exec("DELETE FROM packages")
		db.Exec("DELETE FROM live_classes")
		db.Exec("DELETE FROM videos")
		db.Exec("DELETE FROM modules")
		db.Exec("DELETE FROM vouchers")

		// 2. Insert Official Banners
		banners := []models.Banner{
			{
				ID:        uuid.New(),
				Tag:       "PROGRAM UNGGULAN",
				SubTag:    "Metode KPM",
				Title:     "Matematika Nalaria Realistik (MNR)",
				Subtitle:  "Metode berpikir kritis dan pemecahan masalah matematika kontekstual khas KPM",
				ImageURL:  "https://images.unsplash.com/photo-1509228468518-180dd4864904?w=900&auto=format&fit=crop&q=80",
				Route:     "/package_list",
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ID:        uuid.New(),
				Tag:       "TRYOUT CBT",
				SubTag:    "Skor Otomatis",
				Title:     "Simulasi Ujian & Tryout CBT Online",
				Subtitle:  "Ribuan bank soal berstandar olimpiade dengan evaluasi peringkat nasional",
				ImageURL:  "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=900&auto=format&fit=crop&q=80",
				Route:     "/package_list",
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ID:        uuid.New(),
				Tag:       "KELAS LIVE",
				SubTag:    "Tatap Muka Daring",
				Title:     "Sesi Interaktif Live Class Bersama Master Tutor",
				Subtitle:  "Bimbingan tatap muka daring melalui Zoom & Google Meet berjadwal",
				ImageURL:  "https://images.unsplash.com/photo-1588072432836-e10032774350?w=900&auto=format&fit=crop&q=80",
				Route:     "/live_classes",
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ID:        uuid.New(),
				Tag:       "E-BOOK & MODUL",
				SubTag:    "Unduh PDF",
				Title:     "Modul & Bank Soal Digital format PDF",
				Subtitle:  "Arsip materi belajar, rangkuman rumus, dan pembahasan soal terstruktur",
				ImageURL:  "https://images.unsplash.com/photo-1524178232363-1fb2b075b655?w=900&auto=format&fit=crop&q=80",
				Route:     "/package_list",
				SortOrder: 4,
				IsActive:  true,
			},
		}
		for _, b := range banners {
			db.Create(&b)
		}

		// 3. Insert Official Packages
		p1Disc := 99000.0
		p2Disc := 149000.0
		p3Disc := 119000.0
		p4Disc := 199000.0
		p5Disc := 125000.0

		packages := []models.Package{
			{
				ID:                     uuid.New(),
				Title:                  "Matematika Nalaria Realistik (MNR) SD",
				Description:            "Pembinaan logika bernalar, pemahaman pecahan, deret angka, dan penyelesaian soal cerita aplikatif untuk siswa sekolah dasar.",
				Kelas:                  "Kelas 4-6 SD",
				Jenjang:                "SD",
				Price:                  150000,
				DiscountPrice:          &p1Disc,
				IsDiscountActive:       true,
				MembershipDurationDays: 180,
				Thumbnail:              "https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80",
				IsActive:               true,
			},
			{
				ID:                     uuid.New(),
				Title:                  "Intensif Juara Olimpiade Matematika (OSN / KMNR) SMP",
				Description:            "Bedah tipe soal olimpiade sains dan matematika tingkat kabupaten hingga nasional lengkap dengan trik pemecahan masalah cepat.",
				Kelas:                  "Kelas 7-9 SMP",
				Jenjang:                "SMP",
				Price:                  250000,
				DiscountPrice:          &p2Disc,
				IsDiscountActive:       true,
				MembershipDurationDays: 365,
				Thumbnail:              "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80",
				IsActive:               true,
			},
			{
				ID:                     uuid.New(),
				Title:                  "Mastering Aljabar & Geometri Analitik SMP",
				Description:            "Pendalaman materi SPLDV, fungsi kuadrat, phytagoras, lingkaran, serta bangun ruang dengan metode visual interaktif.",
				Kelas:                  "Kelas 8-9 SMP",
				Jenjang:                "SMP",
				Price:                  175000,
				DiscountPrice:          &p3Disc,
				IsDiscountActive:       true,
				MembershipDurationDays: 180,
				Thumbnail:              "https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=600&auto=format&fit=crop&q=80",
				IsActive:               true,
			},
			{
				ID:                     uuid.New(),
				Title:                  "Persiapan UTBK-SNBT Penalaran Matematika & TPS SMA",
				Description:            "Simulasi tes skolastik, pemahaman bacaan kuantitatif, dan penalaran matematika terstandar seleksi masuk PTN.",
				Kelas:                  "Kelas 10-12 SMA",
				Jenjang:                "SMA",
				Price:                  350000,
				DiscountPrice:          &p4Disc,
				IsDiscountActive:       true,
				MembershipDurationDays: 365,
				Thumbnail:              "https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80",
				IsActive:               true,
			},
			{
				ID:                     uuid.New(),
				Title:                  "Sains & IPA Terpadu Eksploratif SD-SMP",
				Description:            "Pembahasan sains fisika, kimia dasar, dan biologi terintegrasi yang mudah dipahami melalui contoh kontekstual kehidupan sehari-hari.",
				Kelas:                  "Kelas 5-8 SD/SMP",
				Jenjang:                "SMP",
				Price:                  180000,
				DiscountPrice:          &p5Disc,
				IsDiscountActive:       true,
				MembershipDurationDays: 180,
				Thumbnail:              "https://images.unsplash.com/photo-1532094349884-543bc11b234d?w=600&auto=format&fit=crop&q=80",
				IsActive:               true,
			},
		}
		for _, p := range packages {
			db.Create(&p)
		}

		// 4. Insert Official Live Classes
		now := time.Now()
		liveClasses := []models.LiveClass{
			{
				ID:              uuid.New(),
				Title:           "Mastering Trik Penalaran Aljabar & Aritmatika Sosial",
				Description:     "Pembahasan metode bernalar khas KPM dalam memecahkan soal aljabar dan pemodelan matematika cerita.",
				TutorName:       "Dr. Ridwan Hasan Saputra, M.Si",
				TutorRole:       "Founder & Master Coach KPM",
				TutorAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&auto=format&fit=crop&q=80",
				Subject:         "Matematika Nalaria",
				Grade:           "SMP & SMA",
				ScheduledAt:     now.Add(72 * time.Hour),
				DurationMinutes: 90,
				ZoomURL:         "https://zoom.us/j/8921471029?pwd=kpmliveclass",
				IsLive:          false,
				IsActive:        true,
			},
			{
				ID:              uuid.New(),
				Title:           "Strategi Emas Juara Olimpiade Sains & MIPA Nasional",
				Description:     "Bedah pola soal OSN/KSN tingkat kota hingga nasional beserta strategi manajemen waktu pengerjaan ujian CBT.",
				TutorName:       "Kak Aris Fadillah, S.Si",
				TutorRole:       "Koordinator Pembinaan Olimpiade KPM",
				TutorAvatar:     "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&auto=format&fit=crop&q=80",
				Subject:         "Olimpiade MIPA",
				Grade:           "SD & SMP",
				ScheduledAt:     now.Add(96 * time.Hour),
				DurationMinutes: 75,
				ZoomURL:         "https://zoom.us/j/8921471030?pwd=kpmliveclass",
				IsLive:          false,
				IsActive:        true,
			},
			{
				ID:              uuid.New(),
				Title:           "Trik Sukses Menjawab Soal Penalaran Matematika UTBK",
				Description:     "Latihan bedah soal kuantitatif sulit dengan metode eliminasi logika dan penalaran praktis.",
				TutorName:       "Kak Faiz Ramadhan, S.Pd",
				TutorRole:       "Tutor Spesialis Penalaran Skolastik",
				TutorAvatar:     "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&auto=format&fit=crop&q=80",
				Subject:         "Penalaran Kuantitatif",
				Grade:           "SMA / Alumni",
				ScheduledAt:     now.Add(120 * time.Hour),
				DurationMinutes: 90,
				ZoomURL:         "https://zoom.us/j/8921471031?pwd=kpmliveclass",
				IsLive:          false,
				IsActive:        true,
			},
		}
		for _, lc := range liveClasses {
			db.Create(&lc)
		}

		// 5. Insert Official Videos
		videos := []models.Video{
			{
				ID:          uuid.New(),
				Title:       "Pengenalan Konsep Dasar Berpikir Nalaria Realistik (MNR)",
				Description: "Penjelasan filosofi berpikir nalar, mengenali pola masalah matematika, dan cara berpikir sistematis.",
				VideoURL:    "https://www.youtube.com/watch?v=kJQP7kiw5Fk",
				Thumbnail:   "https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80",
				Price:       0,
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Title:       "Trik Berhitung Cepat Operasi Bilangan Pecahan & Desimal",
				Description: "Metode penyelesaian pecahan campuran, perbandingan senilai, dan konversi persen secara cepat tanpa rumus rumit.",
				VideoURL:    "https://www.youtube.com/watch?v=kJQP7kiw5Fk",
				Thumbnail:   "https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=600&auto=format&fit=crop&q=80",
				Price:       0,
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Title:       "Strategi Menyelesaikan Soal Cerita Persamaan Linear (SPLDV)",
				Description: "Langkah-langkah mengubah narasi soal cerita kehidupan sehari-hari menjadi bentuk persamaan matematika yang mudah dipecahkan.",
				VideoURL:    "https://www.youtube.com/watch?v=kJQP7kiw5Fk",
				Thumbnail:   "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80",
				Price:       0,
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Title:       "Kupas Tuntas Soal HOTS Logika Kombinatorika & Peluang",
				Description: "Analisis mendalam konsep permutasi, kombinasi, dan kaidah pencacahan untuk persiapan kompetisi sains dan ujian masuk PTN.",
				VideoURL:    "https://www.youtube.com/watch?v=kJQP7kiw5Fk",
				Thumbnail:   "https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80",
				Price:       0,
				IsActive:    true,
			},
		}
		for _, v := range videos {
			db.Create(&v)
		}

		// 6. Insert Official Modules
		modules := []models.Module{
			{
				ID:          uuid.New(),
				Title:       "Buku Panduan Konsep Matematika Nalaria Realistik (MNR) Dasar",
				Description: "Modul pengantar metode penalaran matematika realistik untuk tingkat SD dan SMP awal.",
				Subject:     "Matematika Nalaria",
				Grade:       "SD & SMP",
				Chapter:     1,
				TotalPages:  36,
				PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				CoverImage:  "https://images.unsplash.com/photo-1509228468518-180dd4864904?w=600&auto=format&fit=crop&q=80",
				IsActive:    true,
				SortOrder:   1,
			},
			{
				ID:          uuid.New(),
				Title:       "Kompilasi Soal Kompetisi Matematika Nalaria (KMNR) & Pembahasan",
				Description: "Kumpulan bank soal kompetisi KMNR babak penyisihan hingga final lengkap dengan pembahasan cara bernalar.",
				Subject:     "Olimpiade MIPA",
				Grade:       "Semua Jenjang",
				Chapter:     2,
				TotalPages:  64,
				PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				CoverImage:  "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=600&auto=format&fit=crop&q=80",
				IsActive:    true,
				SortOrder:   2,
			},
			{
				ID:          uuid.New(),
				Title:       "Rangkuman Rumus Praktis & Pemecahan Soal MIPA SMP-SMA",
				Description: "Panduan saku rumus esensial matematika dan sains terpadu yang dilengkapi contoh soal latihan terarah.",
				Subject:     "Sains Terpadu",
				Grade:       "SMP & SMA",
				Chapter:     3,
				TotalPages:  48,
				PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				CoverImage:  "https://images.unsplash.com/photo-1532094349884-543bc11b234d?w=600&auto=format&fit=crop&q=80",
				IsActive:    true,
				SortOrder:   3,
			},
			{
				ID:          uuid.New(),
				Title:       "Bank Latihan Soal Tryout Penalaran Kuantitatif & TPS",
				Description: "Paket modul simulasi tes potensi skolastik dan penalaran numerik persiapan ujian seleksi perguruan tinggi.",
				Subject:     "Penalaran Kuantitatif",
				Grade:       "SMA / UTBK",
				Chapter:     4,
				TotalPages:  52,
				PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				CoverImage:  "https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=600&auto=format&fit=crop&q=80",
				IsActive:    true,
				SortOrder:   4,
			},
		}
		for _, m := range modules {
			db.Create(&m)
		}

		// 7. Insert Official Vouchers
		vouchers := []models.Voucher{
			{
				ID:              uuid.New(),
				Code:            "KPMJUARA2026",
				Title:           "Diskon Spesial 25% Seluruh Paket",
				Description:     "Potongan harga 25% untuk pembelian semua paket bimbingan dan tryout ujian CBT KPM Academy.",
				DiscountPercent: 25,
				MaxDiscount:     50000,
				ValidUntil:      now.AddDate(1, 0, 0),
				IsActive:        true,
				UsageLimit:      1000,
			},
			{
				ID:             uuid.New(),
				Code:           "MNRINDONESIA",
				Title:          "Potongan Langsung Rp 50.000 Paket MNR",
				Description:    "Spesial paket pembinaan Matematika Nalaria Realistik tingkat SD & SMP.",
				DiscountAmount: 50000,
				MinPurchase:    100000,
				ValidUntil:     now.AddDate(1, 0, 0),
				IsActive:       true,
				UsageLimit:     500,
			},
			{
				ID:              uuid.New(),
				Code:            "PELAJARHEBAT",
				Title:           "Diskon 15% Member Baru KPM",
				Description:     "Voucher potongan 15% untuk seluruh siswa baru KPM Academy.",
				DiscountPercent: 15,
				MaxDiscount:     35000,
				ValidUntil:      now.AddDate(1, 0, 0),
				IsActive:        true,
				UsageLimit:      500,
			},
		}
		for _, v := range vouchers {
			db.Create(&v)
		}

		return utils.SuccessResponse(c, http.StatusOK, "Katalog resmi KPM Academy berhasil dimuat dan ditata ulang!", fiber.Map{
			"total_banners":      len(banners),
			"total_packages":     len(packages),
			"total_live_classes": len(liveClasses),
			"total_videos":       len(videos),
			"total_modules":      len(modules),
			"total_vouchers":     len(vouchers),
		})
	}
}
