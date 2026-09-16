package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Register(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()
		var existing models.User
		if result := db.Where("email = ?", req.Email).First(&existing); result.Error == nil {
			return utils.ConflictResponse(c, "Email already registered")
		}

		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to hash password")
		}

		user := models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hashed,
			Phone:    req.Phone,
			Role:     "user",
		}
		if result := db.Create(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to create user")
		}

		tokens, err := utils.GenerateTokenPair(user.ID, user.Email, user.Role, &cfg.JWT)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to generate tokens")
		}
		user.RefreshToken = tokens.RefreshToken
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to save refresh token")
		}

		return utils.SuccessResponse(c, fiber.StatusCreated, "Registration successful", fiber.Map{
			"user": user, "tokens": tokens,
		})
	}
}

func Login(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()
		var user models.User
		if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
			return utils.UnauthorizedResponse(c, "Invalid email or password")
		}
		if !user.IsActive {
			return utils.UnauthorizedResponse(c, "Account is deactivated")
		}
		if !user.IsVerified {
			return utils.UnauthorizedResponse(c, "Account is not verified. Please verify your email first.")
		}
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			return utils.UnauthorizedResponse(c, "Invalid email or password")
		}

		now := time.Now()
		user.LastLoginAt = &now

		tokens, err := utils.GenerateTokenPair(user.ID, user.Email, user.Role, &cfg.JWT)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to generate tokens")
		}
		user.RefreshToken = tokens.RefreshToken
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to save user")
		}

		log := models.LoginLog{
			UserID:    user.ID,
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			LoginAt:   now,
		}
		if result := db.Create(&log); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to log login")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
			"user": user, "tokens": tokens,
		})
	}
}

func RefreshToken(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RefreshTokenRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		claims, err := utils.ValidateRefreshToken(req.RefreshToken, &cfg.JWT)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid or expired refresh token")
		}

		db := database.GetDB()
		var user models.User
		if result := db.Where("id = ?", claims.UserID).First(&user); result.Error != nil {
			return utils.UnauthorizedResponse(c, "User not found")
		}
		if user.RefreshToken != req.RefreshToken {
			return utils.UnauthorizedResponse(c, "Refresh token mismatch")
		}

		tokens, err := utils.GenerateTokenPair(user.ID, user.Email, user.Role, &cfg.JWT)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to generate tokens")
		}
		user.RefreshToken = tokens.RefreshToken
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to save refresh token")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed", fiber.Map{
			"tokens": tokens,
		})
	}
}

func Logout(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var user models.User
		if result := db.Where("id = ?", uid).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}
		user.RefreshToken = ""
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to logout")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Logged out successfully", nil)
	}
}

func ForgotPassword(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ForgotPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()
		var user models.User
		if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
			return utils.SuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link has been sent", nil)
		}

		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		token := hex.EncodeToString(tokenBytes)

		resetToken := models.PasswordResetToken{
			Email:     req.Email,
			TokenHash: token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
		db.Where("email = ?", req.Email).Delete(&models.PasswordResetToken{})
		db.Create(&resetToken)

		log.Printf("Password reset token for %s: %s (implement email sending in production)", req.Email, token)

		return utils.SuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link has been sent", fiber.Map{
			"reset_token": token,
		})
	}
}

func ResetPassword(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.ResetPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		db := database.GetDB()
		var resetToken models.PasswordResetToken
		if result := db.Where("email = ? AND token_hash = ? AND expires_at > ?",
			req.Email, req.Token, time.Now()).First(&resetToken); result.Error != nil {
			return utils.BadRequestResponse(c, "Invalid or expired reset token", nil)
		}

		hashed, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to hash password")
		}

		var user models.User
		if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "User not found")
		}
		user.Password = hashed
		user.RefreshToken = ""
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to reset password")
		}

		db.Where("email = ?", req.Email).Delete(&models.PasswordResetToken{})

		return utils.SuccessResponse(c, fiber.StatusOK, "Password reset successful", nil)
	}
}

func GetProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()
		var user models.User
		if result := db.Where("id = ?", uid).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved", user)
	}
}

func UpdateProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()

		var user models.User
		if result := db.Where("id = ?", uid).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}

		var req models.UpdateProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}

		if req.Name != "" {
			user.Name = req.Name
		}
		user.Phone = req.Phone
		user.StudentName = req.StudentName
		user.StudentClass = req.StudentClass
		user.StudentMajor = req.StudentMajor
		user.SchoolName = req.SchoolName
		if req.ProfilePhoto != "" {
			user.ProfilePhoto = req.ProfilePhoto
		}
		user.Address = req.Address
		user.Gender = req.Gender
		user.Religion = req.Religion

		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to update profile")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "Profile updated", user)
	}
}

func ChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, ok := utils.GetUserID(c)
		if !ok {
			return utils.UnauthorizedResponse(c, "Unauthorized")
		}
		db := database.GetDB()

		var req models.ChangePasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if errs := utils.ValidateStruct(req); len(errs) > 0 {
			return utils.BadRequestResponse(c, "Validation failed", errs)
		}

		var user models.User
		if result := db.Where("id = ?", uid).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}
		if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
			return utils.BadRequestResponse(c, "Old password is incorrect", nil)
		}

		hashed, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			return utils.InternalErrorResponse(c, "Failed to hash password")
		}
		user.Password = hashed
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to change password")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Password changed", nil)
	}
}

func GetAllUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.User{}).Count(&total)

		var users []models.User
		if result := db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&users); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve users")
		}

		return utils.SuccessResponseWithPagination(c, "Users retrieved", users, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}

func GetUserByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid user ID", nil)
		}
		db := database.GetDB()
		var user models.User
		if result := db.Where("id = ?", id).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "User retrieved", user)
	}
}

func ToggleUserActive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.BadRequestResponse(c, "Invalid user ID", nil)
		}
		db := database.GetDB()
		var user models.User
		if result := db.Where("id = ?", id).First(&user); result.Error != nil {
			return utils.NotFoundResponse(c, "User not found")
		}
		user.IsActive = !user.IsActive
		if result := db.Save(&user); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to toggle user status")
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "User status toggled", user)
	}
}

func GetLoginLogs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		page, perPage, offset := utils.GetPagination(c)

		var total int64
		db.Model(&models.LoginLog{}).Count(&total)

		var logs []models.LoginLog
		if result := db.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&logs); result.Error != nil {
			return utils.InternalErrorResponse(c, "Failed to retrieve login logs")
		}

		return utils.SuccessResponseWithPagination(c, "Login logs retrieved", logs, &utils.Pagination{
			Total: total, PerPage: perPage, CurrentPage: page, LastPage: utils.CalcLastPage(total, perPage),
		})
	}
}
