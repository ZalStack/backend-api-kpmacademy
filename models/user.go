package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	StudentName     string         `gorm:"type:varchar(255)" json:"student_name"`
	StudentClass    string         `gorm:"type:varchar(50)" json:"student_class"`
	StudentMajor    string         `gorm:"type:varchar(100)" json:"student_major"`
	SchoolName      string         `gorm:"type:varchar(255)" json:"school_name"`
	ProfilePhoto    string         `gorm:"type:varchar(255)" json:"profile_photo"`
	Address         string         `gorm:"type:text" json:"address"`
	Gender          string         `gorm:"type:varchar(20)" json:"gender"`
	Religion        string         `gorm:"type:varchar(20)" json:"religion"`
	Role            string         `gorm:"type:varchar(50);default:user;not null" json:"role"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	RefreshToken    string         `gorm:"type:text" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	StudentName  string `json:"student_name"`
	StudentClass string `json:"student_class"`
	StudentMajor string `json:"student_major"`
	SchoolName   string `json:"school_name"`
	ProfilePhoto string `json:"profile_photo"`
	Address      string `json:"address"`
	Gender       string `json:"gender"`
	Religion     string `json:"religion"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}
