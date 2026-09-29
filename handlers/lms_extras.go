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

// GetLeaderboard returns ranked students based on practice sessions and scores
func GetLeaderboard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		period := c.Query("period", "weekly")
		limit := c.QueryInt("limit", 50)
		if limit <= 0 || limit > 100 {
			limit = 50
		}

		var dateFilter time.Time
		now := time.Now()
		switch period {
		case "weekly":
			dateFilter = now.AddDate(0, 0, -7)
		case "monthly":
			dateFilter = now.AddDate(0, -1, 0)
		default:
			dateFilter = time.Time{}
		}

		type QueryResult struct {
			UserID        string  `json:"user_id"`
			Name          string  `json:"name"`
			StudentClass  string  `json:"student_class"`
			SchoolName    string  `json:"school_name"`
			ProfilePhoto  string  `json:"profile_photo"`
			AvgScore      float64 `json:"avg_score"`
			TotalPractice int     `json:"total_practice"`
			TotalCorrect  int     `json:"total_correct"`
		}

		var results []QueryResult
		q := db.Table("practice_sessions ps").
			Select(`
				u.id as user_id,
				u.name,
				u.student_class,
				u.school_name,
				u.profile_photo,
				COALESCE(AVG(ps.total_score), 0) as avg_score,
				COUNT(ps.id) as total_practice,
				COALESCE(SUM(ps.correct_answer), 0) as total_correct
			`).
			Joins("JOIN users u ON u.id = ps.user_id").
			Where("ps.status = ?", "finished")

		if !dateFilter.IsZero() {
			q = q.Where("ps.created_at >= ?", dateFilter)
		}

		_ = q.Group("u.id, u.name, u.student_class, u.school_name, u.profile_photo").
			Order("avg_score DESC, total_correct DESC").
			Limit(limit).
			Scan(&results).Error

		var entries []models.LeaderboardResponseItem
		for i, r := range results {
			avatar := r.ProfilePhoto
			if avatar == "" {
				avatar = "https://api.dicebear.com/7.x/adventurer/png?seed=" + r.Name
			}
			entries = append(entries, models.LeaderboardResponseItem{
				Rank:          i + 1,
				UserID:        r.UserID,
				Name:          r.Name,
				StudentClass:  r.StudentClass,
				SchoolName:    r.SchoolName,
				Avatar:        avatar,
				AvgScore:      r.AvgScore,
				TotalPractice: r.TotalPractice,
				TotalCorrect:  r.TotalCorrect,
			})
		}

		// Fallback sample leaderboard if database has fewer than 3 records
		if len(entries) < 3 {
			entries = []models.LeaderboardResponseItem{
				{
					Rank:          1,
					UserID:        "usr_top_1",
					Name:          "Ahmad Faiz Al-Ghifari",
					StudentClass:  "XII IPA 1",
					SchoolName:    "SMAN 1 Surakarta (Target STEI ITB)",
					Avatar:        "https://api.dicebear.com/7.x/adventurer/png?seed=Faiz",
					AvgScore:      985.0,
					TotalPractice: 24,
					TotalCorrect:  472,
				},
				{
					Rank:          2,
					UserID:        "usr_top_2",
					Name:          "Nadia Salsabila Putri",
					StudentClass:  "XII MIPA",
					SchoolName:    "SMA Kharisma Bangsa (Target FK UI)",
					Avatar:        "https://api.dicebear.com/7.x/adventurer/png?seed=Nadia",
					AvgScore:      970.5,
					TotalPractice: 21,
					TotalCorrect:  410,
				},
				{
					Rank:          3,
					UserID:        "usr_top_3",
					Name:          "Michael Kevin Wijaya",
					StudentClass:  "XI SMA",
					SchoolName:    "SMA Santa Angela (Olimpiade MNR)",
					Avatar:        "https://api.dicebear.com/7.x/adventurer/png?seed=Kevin",
					AvgScore:      955.0,
					TotalPractice: 19,
					TotalCorrect:  375,
				},
				{
					Rank:          4,
					UserID:        "usr_top_4",
					Name:          "Zahra Annisa Rahma",
					StudentClass:  "XII IPA 3",
					SchoolName:    "SMAN 8 Jakarta (Target FTI ITB)",
					Avatar:        "https://api.dicebear.com/7.x/adventurer/png?seed=Zahra",
					AvgScore:      940.0,
					TotalPractice: 18,
					TotalCorrect:  350,
				},
				{
					Rank:          5,
					UserID:        "usr_top_5",
					Name:          "Rizky Pratama Yudha",
					StudentClass:  "XII IPA",
					SchoolName:    "SMAN 3 Bandung (Target FK Unpad)",
					Avatar:        "https://api.dicebear.com/7.x/adventurer/png?seed=Rizky",
					AvgScore:      925.5,
					TotalPractice: 16,
					TotalCorrect:  312,
				},
			}
		}

		return utils.SuccessResponse(c, http.StatusOK, "Papan skor nasional berhasil dimuat", entries)
	}
}

// GetAllTutors returns active tutors with their specialty tags and available slots
func GetAllTutors() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		var tutors []models.Tutor
		db.Where("is_active = ?", true).Order("sort_order asc, created_at desc").Find(&tutors)

		if len(tutors) == 0 {
			seedTutors := []models.Tutor{
				{
					ID:              uuid.New(),
					Name:            "Dr. Ir. R. Ridwan Hasan Saputra, M.Si.",
					Role:            "Master Coach Matematika Nalaria (MNR)",
					Experience:      "25+ Thn Pengalaman",
					Rating:          5.0,
					TotalSessions:   3500,
					Avatar:          "https://api.dicebear.com/7.x/adventurer/png?seed=Ridwan",
					PricePerSession: 150000,
					Specialties:     "MNR,Olimpiade,Pelatihan Guru",
					AvailableSlots:  "09:00 - 10:30,13:30 - 15:00,19:30 - 21:00",
					IsActive:        true,
					SortOrder:       1,
				},
				{
					ID:              uuid.New(),
					Name:            "Budi Pratama, M.Sc.",
					Role:            "Master Tutor Fisika & Sains Terapan",
					Experience:      "10+ Thn Pengalaman",
					Rating:          4.9,
					TotalSessions:   1420,
					Avatar:          "https://api.dicebear.com/7.x/adventurer/png?seed=Budi",
					PricePerSession: 120000,
					Specialties:     "Fisika SMA,Olimpiade Sains",
					AvailableSlots:  "10:00 - 11:30,16:00 - 17:30,20:00 - 21:30",
					IsActive:        true,
					SortOrder:       2,
				},
				{
					ID:              uuid.New(),
					Name:            "Siti Rahma, S.Si.",
					Role:            "Spesialis Matematika SD & SMP",
					Experience:      "8+ Thn Pengalaman",
					Rating:          4.9,
					TotalSessions:   1150,
					Avatar:          "https://api.dicebear.com/7.x/adventurer/png?seed=Siti",
					PricePerSession: 100000,
					Specialties:     "Matematika SD,Logika Berhitung",
					AvailableSlots:  "08:30 - 10:00,14:00 - 15:30,16:30 - 18:00",
					IsActive:        true,
					SortOrder:       3,
				},
				{
					ID:              uuid.New(),
					Name:            "Ahmad Fauzi, M.Pd.",
					Role:            "Tutor Kimia & Asesmen Sains",
					Experience:      "7+ Thn Pengalaman",
					Rating:          4.8,
					TotalSessions:   890,
					Avatar:          "https://api.dicebear.com/7.x/adventurer/png?seed=Fauzi",
					PricePerSession: 100000,
					Specialties:     "Kimia Terapan,Eksperimen MIPA",
					AvailableSlots:  "11:00 - 12:30,15:30 - 17:00,18:30 - 20:00",
					IsActive:        true,
					SortOrder:       4,
				},
			}
			for _, t := range seedTutors {
				db.Create(&t)
			}
			tutors = seedTutors
		}

		type TutorResponse struct {
			ID              string   `json:"id"`
			Name            string   `json:"name"`
			Role            string   `json:"role"`
			Experience      string   `json:"experience"`
			Rating          float64  `json:"rating"`
			TotalSessions   int      `json:"total_sessions"`
			Avatar          string   `json:"avatar"`
			PricePerSession float64  `json:"price_per_session"`
			Specialties     []string `json:"specialties"`
			AvailableSlots  []string `json:"available_slots"`
		}

		var resp []TutorResponse
		for _, t := range tutors {
			var specs []string
			if t.Specialties != "" {
				specs = strings.Split(t.Specialties, ",")
			}
			var slots []string
			if t.AvailableSlots != "" {
				slots = strings.Split(t.AvailableSlots, ",")
			}
			resp = append(resp, TutorResponse{
				ID:              t.ID.String(),
				Name:            t.Name,
				Role:            t.Role,
				Experience:      t.Experience,
				Rating:          t.Rating,
				TotalSessions:   t.TotalSessions,
				Avatar:          t.Avatar,
				PricePerSession: t.PricePerSession,
				Specialties:     specs,
				AvailableSlots:  slots,
			})
		}

		return utils.SuccessResponse(c, http.StatusOK, "Daftar tutor berhasil dimuat", resp)
	}
}

// BookTutorSchedule handles 1-on-1 private mentoring session booking
func BookTutorSchedule() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		user := c.Locals("user").(*models.User)

		tutorParam := c.Params("id")
		tutorID, err := uuid.Parse(tutorParam)
		if err != nil {
			return utils.BadRequestResponse(c, "ID Tutor tidak valid", nil)
		}

		var tutor models.Tutor
		if err := db.Where("id = ? AND is_active = ?", tutorID, true).First(&tutor).Error; err != nil {
			return utils.NotFoundResponse(c, "Tutor tidak ditemukan")
		}

		var req models.BookTutorRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Payload reservasi tidak valid", err.Error())
		}

		if req.ScheduledDate == "" || req.ScheduledTimeSlot == "" || req.Subject == "" {
			return utils.BadRequestResponse(c, "Tanggal, waktu sesi, dan topik bimbingan wajib diisi", nil)
		}

		booking := models.TutorBooking{
			ID:                uuid.New(),
			TutorID:           tutor.ID,
			TutorName:         tutor.Name,
			StudentID:         user.ID,
			StudentName:       user.Name,
			ScheduledDate:     req.ScheduledDate,
			ScheduledTimeSlot: req.ScheduledTimeSlot,
			Subject:           req.Subject,
			Status:            "confirmed",
			MeetingURL:        "https://zoom.us/j/kpm_tutor_" + tutor.ID.String()[:8],
			Notes:             req.Notes,
		}

		if err := db.Create(&booking).Error; err != nil {
			return utils.InternalServerError(c, "Gagal membuat reservasi tutor: "+err.Error())
		}

		return utils.SuccessResponse(c, http.StatusCreated, "Jadwal les privat 1-on-1 berhasil dipesan!", booking)
	}
}

// GetMyTutorBookings returns all private sessions booked by the authenticated student
func GetMyTutorBookings() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		user := c.Locals("user").(*models.User)

		var bookings []models.TutorBooking
		db.Where("student_id = ?", user.ID).Order("created_at desc").Find(&bookings)

		return utils.SuccessResponse(c, http.StatusOK, "Riwayat booking tutor berhasil dimuat", bookings)
	}
}

