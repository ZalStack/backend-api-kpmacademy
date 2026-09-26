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

// GetAllBanners returns active banners sorted by sort_order
func GetAllBanners() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var banners []models.Banner
		db.Where("is_active = ?", true).Order("sort_order asc, created_at desc").Find(&banners)

		// Seed initial banners if table is empty
		if len(banners) == 0 {
			seedBanners := []models.Banner{
				{
					ID:        uuid.New(),
					Tag:       "VIDEO MATERI",
					SubTag:    "Akses Fleksibel",
					Title:     "Video Pembelajaran Interaktif",
					Subtitle:  "Pelajari konsep materi dari dasar hingga mahir kapan saja dan di mana saja",
					ImageURL:  "https://images.unsplash.com/photo-1524178232363-1fb2b075b655?w=900&auto=format&fit=crop&q=80",
					Route:     "/video_list",
					IsActive:  true,
					SortOrder: 1,
				},
				{
					ID:        uuid.New(),
					Tag:       "TRYOUT CBT",
					SubTag:    "Sistem Skor IRT",
					Title:     "Simulasi Ujian & Tryout CBT",
					Subtitle:  "Asah kemampuan dengan ribuan bank soal dan pembahasan lengkap",
					ImageURL:  "https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=900&auto=format&fit=crop&q=80",
					Route:     "/package_list",
					IsActive:  true,
					SortOrder: 2,
				},
				{
					ID:        uuid.New(),
					Tag:       "PAKET UNGGULAN",
					SubTag:    "Diskon 50%",
					Title:     "Paket Intensif MIPA KPM",
					Subtitle:  "Kuasai konsep Matematika Nalaria & Sains bersama Master Tutor",
					ImageURL:  "https://images.unsplash.com/photo-1509228468518-180dd4864904?w=900&auto=format&fit=crop&q=80",
					Route:     "/package_list",
					IsActive:  true,
					SortOrder: 3,
				},
			}
			for _, b := range seedBanners {
				db.Create(&b)
			}
			banners = seedBanners
		}

		return utils.SuccessResponse(c, http.StatusOK, "Banners retrieved successfully", banners)
	}
}

// GetAllLiveClasses returns active Live Classes with zoom link and tutor details
func GetAllLiveClasses() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var liveClasses []models.LiveClass
		db.Where("is_active = ?", true).Order("scheduled_at asc").Find(&liveClasses)

		if len(liveClasses) == 0 {
			now := time.Now()
			seedLiveClasses := []models.LiveClass{
				{
					ID:              uuid.New(),
					Title:           "Mastering Trik Penalaran Aljabar & Aritmatika Sosial",
					Description:     "Pembahasan metode cepat Nalaria Realistik dalam memecahkan soal aljabar dan pemodelan matematika cerita.",
					TutorName:       "Dr. Ridwan Hasan Saputra, M.Si",
					TutorRole:       "Founder & Master Coach KPM",
					TutorAvatar:     "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&auto=format&fit=crop&q=80",
					Subject:         "Matematika Nalaria",
					Grade:           "SMP & SMA",
					ScheduledAt:     now.Add(2 * time.Hour),
					DurationMinutes: 90,
					ZoomURL:         "https://zoom.us/j/8921471029?pwd=kpmliveclassdemo",
					IsLive:          true,
					IsActive:        true,
				},
				{
					ID:              uuid.New(),
					Title:           "Strategi Emas Juara Olimpiade Sains & MIPA Nasional",
					Description:     "Bedah pola soal OSN/KSN tingkat kota hingga nasional beserta strategi manajemen waktu ujian CBT.",
					TutorName:       "Kak Aris Fadillah, S.Si",
					TutorRole:       "Koordinator Pembinaan Olimpiade KPM",
					TutorAvatar:     "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&auto=format&fit=crop&q=80",
					Subject:         "Sains & MIPA",
					Grade:           "SD & SMP",
					ScheduledAt:     now.Add(24 * time.Hour),
					DurationMinutes: 75,
					ZoomURL:         "https://zoom.us/j/8921471030?pwd=kpmliveclassdemo",
					IsLive:          false,
					IsActive:        true,
				},
			}
			for _, lc := range seedLiveClasses {
				db.Create(&lc)
			}
			liveClasses = seedLiveClasses
		}

		return utils.SuccessResponse(c, http.StatusOK, "Live classes retrieved successfully", liveClasses)
	}
}

// GetAllModules returns PDF reading modules and learning materials
func GetAllModules() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		subject := c.Query("subject")
		grade := c.Query("grade")

		query := db.Where("is_active = ?", true)
		if subject != "" && subject != "Semua" {
			query = query.Where("LOWER(subject) = ?", strings.ToLower(subject))
		}
		if grade != "" && grade != "Semua" {
			query = query.Where("LOWER(grade) = ?", strings.ToLower(grade))
		}

		var modules []models.Module
		query.Order("chapter asc, sort_order asc").Find(&modules)

		if len(modules) == 0 {
			seedModules := []models.Module{
				{
					ID:          uuid.New(),
					Title:       "Konsep Dasar & Filosofi Matematika Nalaria Realistik (MNR)",
					Description: "Panduan pengenalan metode bernalar khas KPM untuk melatih logika berpikir runtut dan sistematis.",
					Subject:     "Matematika Nalaria",
					Grade:       "SD & SMP",
					Chapter:     1,
					TotalPages:  24,
					PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
					IsActive:    true,
					SortOrder:   1,
				},
				{
					ID:          uuid.New(),
					Title:       "Trik Berhitung Cepat & Operasi Bilangan Bulat Pecahan",
					Description: "Kumpulan trik perhitungan cepat operasi pecahan, persen, dan desimal tanpa rumus rumit.",
					Subject:     "Matematika Nalaria",
					Grade:       "SD & SMP",
					Chapter:     2,
					TotalPages:  36,
					PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
					IsActive:    true,
					SortOrder:   2,
				},
				{
					ID:          uuid.New(),
					Title:       "Bedah Pola Soal Cerita & Pemodelan Aljabar Kontekstual",
					Description: "Strategi mengubah permasalahan sehari-hari menjadi kalimat matematika secara cepat dan akurat.",
					Subject:     "Aljabar & Logika",
					Grade:       "SMP & SMA",
					Chapter:     3,
					TotalPages:  48,
					PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
					IsActive:    true,
					SortOrder:   3,
				},
				{
					ID:          uuid.New(),
					Title:       "Kompilasi Soal Kompetisi Sains & Pembahasan Olimpiade",
					Description: "Bank soal olimpiade sains terpadu lengkap dengan kunci dan teknik eliminasi jawaban analitis.",
					Subject:     "Olimpiade MIPA",
					Grade:       "Semua Jenjang",
					Chapter:     4,
					TotalPages:  52,
					PDFURL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
					IsActive:    true,
					SortOrder:   4,
				},
			}
			for _, m := range seedModules {
				db.Create(&m)
			}
			modules = seedModules
		}

		return utils.SuccessResponse(c, http.StatusOK, "Modules retrieved successfully", modules)
	}
}

// ClaimVoucher validates and claims a promo discount voucher
func ClaimVoucher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ClaimVoucherRequest
		if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Code) == "" {
			return utils.BadRequestResponse(c, "Kode voucher harus diisi", nil)
		}

		code := strings.ToUpper(strings.TrimSpace(req.Code))
		db := database.GetDB()

		var voucher models.Voucher
		result := db.Where("code = ? AND is_active = ?", code, true).First(&voucher)
		if result.Error != nil {
			// Check standard promo vouchers
			switch code {
			case "KPMJUARA2026":
				voucher = models.Voucher{
					ID:              uuid.New(),
					Code:            "KPMJUARA2026",
					Title:           "Diskon Spesial 25% Semua Paket",
					Description:     "Potongan langsung 25% untuk seluruh paket belajar dan persiapan tryout CBT.",
					DiscountPercent: 25,
					MaxDiscount:     50000,
					ValidUntil:      time.Now().AddDate(0, 1, 0),
					IsActive:        true,
				}
				db.Create(&voucher)
			case "MNRINDONESIA":
				voucher = models.Voucher{
					ID:              uuid.New(),
					Code:            "MNRINDONESIA",
					Title:           "Potongan Rp 50.000 Paket MNR",
					Description:     "Spesial pembinaan Matematika Nalaria Realistik tingkat SD & SMP.",
					DiscountAmount:  50000,
					MinPurchase:     100000,
					ValidUntil:      time.Now().AddDate(0, 1, 0),
					IsActive:        true,
				}
				db.Create(&voucher)
			case "MIPABERSAMA":
				voucher = models.Voucher{
					ID:              uuid.New(),
					Code:            "MIPABERSAMA",
					Title:           "Cashback Rp 20.000 Belajar Sains",
					Description:     "Voucher potongan belajar grup sains dan olimpiade MIPA nasional.",
					DiscountAmount:  20000,
					MinPurchase:     50000,
					ValidUntil:      time.Now().AddDate(0, 1, 0),
					IsActive:        true,
				}
				db.Create(&voucher)
			default:
				return utils.NotFoundResponse(c, "Kode voucher tidak ditemukan atau sudah tidak berlaku.")
			}
		}

		if !voucher.ValidUntil.IsZero() && time.Now().After(voucher.ValidUntil) {
			return utils.BadRequestResponse(c, "Masa berlaku voucher telah berakhir", nil)
		}

		return utils.SuccessResponse(c, http.StatusOK, "Voucher berhasil diklaim!", fiber.Map{
			"voucher": voucher,
			"code":    voucher.Code,
			"title":   voucher.Title,
			"message": "Kode voucher aktif dan siap digunakan saat checkout paket!",
		})
	}
}
