package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Testimonial struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Rating     int            `gorm:"default:5" json:"rating"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	IsApproved bool           `gorm:"default:false" json:"is_approved"`
	ApprovedAt *time.Time     `json:"approved_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Testimonial) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type CreateTestimonialRequest struct {
	Content string `json:"content" validate:"required"`
	Rating  int    `json:"rating"`
}

type ApproveTestimonialRequest struct {
	IsApproved bool `json:"is_approved"`
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids" validate:"required"`
}
