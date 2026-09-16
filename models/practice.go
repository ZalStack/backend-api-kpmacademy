package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PracticeSession struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PackageID      uuid.UUID      `gorm:"type:char(36);not null" json:"package_id"`
	Package        Package        `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	OrderID        uuid.UUID      `gorm:"type:char(36);not null" json:"order_id"`
	Order          Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CardID         string         `gorm:"type:varchar(100);not null" json:"card_id"`
	TotalQuestion  int            `gorm:"default:0" json:"total_question"`
	CorrectAnswer  int            `gorm:"default:0" json:"correct_answer"`
	WrongAnswer    int            `gorm:"default:0" json:"wrong_answer"`
	Unanswered     int            `gorm:"default:0" json:"unanswered"`
	TotalScore     float64        `gorm:"type:decimal(5,2);default:0" json:"total_score"`
	DurationSeconds int           `gorm:"default:0" json:"duration_seconds"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
	Status         string         `gorm:"type:varchar(50);default:in_progress" json:"status"`
	Answers        string         `gorm:"type:json" json:"answers"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (p *PracticeSession) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type StartPracticeRequest struct {
	PackageID string `json:"package_id" validate:"required"`
	CardID    string `json:"card_id" validate:"required"`
}

type SubmitPracticeRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Answers   string `json:"answers" validate:"required"`
	Duration  int    `json:"duration"`
}
