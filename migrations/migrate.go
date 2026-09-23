package migrations

import (
	"errors"
	"log"
	"os"
	"strings"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RunMigrations() {
	db := database.GetDB()

	err := db.AutoMigrate(
		&models.User{},
		&models.Package{},
		&models.Order{},
		&models.PracticeSession{},
		&models.SupportTicket{},
		&models.Testimonial{},
		&models.Video{},
		&models.VideoOrder{},
		&models.LoginLog{},
		&models.Notification{},
		&models.ContactForm{},
		&models.ChatMessage{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func SeedAdmin() {
	db := database.GetDB()

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Fatal("FATAL: ADMIN_EMAIL environment variable is required for admin seed")
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("FATAL: ADMIN_PASSWORD environment variable is required for admin seed")
	}

	// Idempotent check by email (including soft-deleted rows)
	// Menggunakan Unscoped agar tidak terkecoh soft-delete: unique index MySQL tetap blokir email yang soft-deleted
	var existing models.User
	err := db.Unscoped().Where("email = ?", adminEmail).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"role":        "admin",
			"is_verified": true,
			"is_active":   true,
		}
		if existing.DeletedAt.Valid {
			updates["deleted_at"] = nil
		}
		if !utils.CheckPasswordHash(adminPassword, existing.Password) {
			newHash, hashErr := utils.HashPassword(adminPassword)
			if hashErr == nil {
				updates["password"] = newHash
			}
		}
		if err := db.Unscoped().Model(&existing).Updates(updates).Error; err != nil {
			log.Printf("Failed to update existing admin: %v", err)
			return
		}
		log.Println("Admin verified and updated successfully")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to check existing admin: %v", err)
		return
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := models.User{
		Name:       "Admin KPM Academy",
		Email:      adminEmail,
		Password:   hashedPassword,
		Role:       "admin",
		IsVerified: true,
		IsActive:   true,
	}

	// OnConflict DoNothing untuk handle race condition / duplicate key (Error 1062)
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&admin)
	if result.Error != nil {
		// Fallback: jika driver tidak support DoNothing, cek duplicate secara string
		if strings.Contains(result.Error.Error(), "Duplicate") || strings.Contains(result.Error.Error(), "duplicate") {
			log.Println("Admin already exists, skipping seed")
			return
		}
		log.Printf("Failed to seed admin: %v", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		log.Println("Admin already exists, skipping seed")
		return
	}

	log.Println("Admin seeded successfully")
}
