package migrations

import (
	"log"
	"os"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"
	"backend-api-kpmacademy/utils"
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

func SeedAdmin(cfg interface{}) {
	db := database.GetDB()

	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		log.Println("Admin already exists, skipping seed")
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@kpmacademy.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin123!"
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

	if result := db.Create(&admin); result.Error != nil {
		log.Printf("Failed to seed admin: %v", result.Error)
		return
	}

	log.Println("Admin seeded successfully")
}
