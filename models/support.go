package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupportTicket struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID   string         `gorm:"type:varchar(100)" json:"session_id"`
	UserID      *uuid.UUID     `gorm:"type:char(36)" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string         `gorm:"type:varchar(255)" json:"name"`
	Email       string         `gorm:"type:varchar(255)" json:"email"`
	Question    string         `gorm:"type:text;not null" json:"question"`
	Answer      string         `gorm:"type:text" json:"answer"`
	Status      string         `gorm:"type:enum('pending','answered','closed');default:pending" json:"status"`
	AnsweredAt  *time.Time     `json:"answered_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *SupportTicket) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type SubmitSupportRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Question string `json:"question" validate:"required"`
}

type AnswerSupportRequest struct {
	Answer string `json:"answer" validate:"required"`
}

type UpdateSupportStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
