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
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.ContactForm{}).Where("user_id = ?", uid).Count(&total)

		var forms []models.ContactForm
		if result := db.Where("user_id = ?", uid).Offset(offset).Limit(perPage).Order("created_at DESC").Find(&forms); result.Error != nil {
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
		if uid, ok := utils.GetUserID(c); ok {
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
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid form ID", nil)
		}
		db := database.GetDB()
		var form models.ContactForm
		if result := db.Where("id = ? AND user_id = ?", id, uid).First(&form); result.Error != nil {
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
		if form.ReadAt == nil && form.Status == "pending" {
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
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		now := time.Now()
		form.AdminReply = req.AdminReply
		form.Status = "replied"
		form.RepliedAt = &now
		if result := db.Save(&form); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to submit reply")
		}

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
		if !utils.ValidateStatus(req.Status, []string{"pending", "read", "replied", "closed"}) {
			return utils.BadRequestResponse(c, "Invalid status. Must be: pending, read, replied, or closed", nil)
		}

		form.Status = req.Status
		if result := db.Save(&form); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to update status")
		}
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
		if result := db.Delete(&form); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to delete contact form")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Contact form deleted", nil)
	}
}

func generateAITutorResponse(query string) string {
	q := strings.ToLower(query)
	if strings.Contains(q, "logaritma") {
		return "💡 **Konsep & Trik Cepat Logaritma:**\n\n1️⃣ **Sifat Dasar:**\n• ^a\\log(b × c) = ^a\\log(b) + ^a\\log(c)\n• ^a\\log(b / c) = ^a\\log(b) - ^a\\log(c)\n• ^a\\log(b^n) = n × ^a\\log(b)\n\n2️⃣ **Trik Cepat:**\nJika ^a\\log(b) = x, maka a^x = b. Samakan basis numerus terlebih dahulu sebelum mengoperasikan penjumlahan atau pengurangan!\n\nAda contoh soal logaritma yang ingin kamu bahas bersama?"
	}
	if strings.Contains(q, "stoikiometri") || strings.Contains(q, "kimia") {
		return "🧪 **Panduan Stoikiometri Kimia:**\n\n1️⃣ **Rumus Utama Mol (n):**\n• n = massa (gram) / Mr\n• n = Volume gas (STP) / 22.4 Liter\n• n = Molaritas (M) × Volume (L)\n\n2️⃣ **Langkah Penyelesaian:**\n1. Setarakan persamaan reaksi kimia.\n2. Ubah semua besaran yang diketahui ke satuan **Mol**.\n3. Gunakan perbandingan koefisien untuk mencari mol zat yang ditanyakan.\n4. Konversi mol tersebut ke gram atau liter sesuai pertanyaan!\n\nSilakan tuliskan soal kimiamu untuk kita bahas langkah demi langkah."
	}
	if strings.Contains(q, "kuadrat") || strings.Contains(q, "aljabar") {
		return "📐 **Rumus Kuadrat Sempurna & Aljabar:**\n\n1️⃣ Bentuk: (a + b)² = a² + 2ab + b²\n2️⃣ Bentuk: (a - b)² = a² - 2ab + b²\n3️⃣ Selisih Kuadrat: a² - b² = (a + b)(a - b)\n\n⚡ **Rumus ABC untuk akar persamaan kuadrat ax² + bx + c = 0:**\nx = (-b ± √(b² - 4ac)) / (2a)\n\nIngin membahas variasi soal kuadrat atau pemfaktoran lainnya?"
	}
	if strings.Contains(q, "mnr") || strings.Contains(q, "nalaria") || strings.Contains(q, "matematika") {
		return "✨ **Metode Matematika Nalaria Realistik (MNR) KPM:**\n\nMetode MNR melatih 3 tahapan berpikir nalar:\n1. **Memahami Masalah:** Identifikasi fakta, pola angka, dan kondisi soal.\n2. **Bernalar Kreatif:** Mengubah soal cerita kompleks menjadi model logika sederhana tanpa rumus kaku.\n3. **Verifikasi:** Mengecek kembali kewajaran hasil jawaban.\n\nKetikkan soal matematika MNR yang sedang kamu pelajari, AI Tutor siap membimbing pengerjaannya!"
	}
	return fmt.Sprintf("Halo! Pertanyaan bagus mengenai \"%s\". 🤖✨\n\n🧠 **Langkah & Solusi Pembahasan:**\n1️⃣ **Identifikasi Konsep:** Analisis variabel atau teori dasar yang terkait dengan topik ini.\n2️⃣ **Langkah Pengerjaan:** Terapkan prinsip step-by-step untuk memecahkan persoalan secara runtut dan tepat.\n3️⃣ **Tips Kunci:** Periksa kembali kesesuaian unit, formula, dan logika penalaran.\n\nAda bagian spesifik atau contoh soal yang ingin kita bedah bersama? Silakan ketikkan ya!", query)
}

func SendChatMessage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()

		var req models.SendChatRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		var userID *uuid.UUID
		if uid, ok := utils.GetUserID(c); ok {
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
		if result := db.Create(&userMsg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to save message")
		}

		aiReplyText, err := utils.CallGeminiAI(req.Message)
		if err != nil || aiReplyText == "" {
			aiReplyText = generateAITutorResponse(req.Message)
		}

		assistantMsg := models.ChatMessage{
			SessionID: sessionID,
			UserID:    userID,
			Role:      "assistant",
			Message:   aiReplyText,
			IsAI:      true,
		}
		if result := db.Create(&assistantMsg); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to save reply")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Message sent", fiber.Map{
			"session_id":        sessionID,
			"reply":             aiReplyText,
			"user_message":      userMsg,
			"assistant_message": assistantMsg,
		})
	}
}

func GetChatHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		sessionID := c.Query("session_id")
		if sessionID == "" {
			return utils.BadRequestResponse(c, "session_id is required", nil)
		}
		db := database.GetDB()

		var messages []models.ChatMessage
		if result := db.Where("session_id = ? AND user_id = ?", sessionID, uid).Order("created_at ASC").Find(&messages); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve chat history")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Chat history retrieved", messages)
	}
}
