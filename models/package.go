package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Package struct {
	ID                   uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title                string         `gorm:"type:varchar(255);not null" json:"title"`
	Description          string         `gorm:"type:text" json:"description"`
	Thumbnail            string         `gorm:"type:varchar(255)" json:"thumbnail"`
	Kelas                string         `gorm:"type:varchar(50)" json:"kelas"`
	Jenjang              string         `gorm:"type:varchar(50)" json:"jenjang"`
	Price                float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	DiscountPrice        *float64       `gorm:"type:decimal(10,2)" json:"discount_price"`
	IsDiscountActive     bool           `gorm:"default:false" json:"is_discount_active"`
	IsPayWhatYouWant     bool           `gorm:"default:false" json:"is_pay_what_you_want"`
	MinPayAmount         float64        `gorm:"type:decimal(10,2);default:0" json:"min_pay_amount"`
	MembershipDurationDays int          `gorm:"default:30" json:"membership_duration_days"`
	Cards                string         `gorm:"type:json" json:"cards"`
	Questions            string         `gorm:"type:json" json:"questions"`
	Reviews              string         `gorm:"type:json" json:"reviews"`
	HideExplanation      bool           `gorm:"default:false" json:"hide_explanation"`
	TimeLimitMinutes     *int           `json:"time_limit_minutes"`
	IsActive             bool           `gorm:"default:true" json:"is_active"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Package) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type CreatePackageRequest struct {
	Title                  string  `json:"title" validate:"required"`
	Description            string  `json:"description"`
	Thumbnail              string  `json:"thumbnail"`
	Kelas                  string  `json:"kelas"`
	Jenjang                string  `json:"jenjang"`
	Price                  float64 `json:"price"`
	DiscountPrice          float64 `json:"discount_price"`
	IsDiscountActive       bool    `json:"is_discount_active"`
	IsPayWhatYouWant       bool    `json:"is_pay_what_you_want"`
	MinPayAmount           float64 `json:"min_pay_amount"`
	MembershipDurationDays int     `json:"membership_duration_days"`
	Cards                  string  `json:"cards"`
	Questions              string  `json:"questions"`
	HideExplanation        bool    `json:"hide_explanation"`
	TimeLimitMinutes       int     `json:"time_limit_minutes"`
}

type UpdatePackageRequest struct {
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	Thumbnail              string   `json:"thumbnail"`
	Kelas                  string   `json:"kelas"`
	Jenjang                string   `json:"jenjang"`
	Price                  *float64 `json:"price"`
	DiscountPrice          *float64 `json:"discount_price"`
	IsDiscountActive       *bool    `json:"is_discount_active"`
	IsPayWhatYouWant       *bool    `json:"is_pay_what_you_want"`
	MinPayAmount           *float64 `json:"min_pay_amount"`
	MembershipDurationDays *int     `json:"membership_duration_days"`
	Cards                  string   `json:"cards"`
	Questions              string   `json:"questions"`
	HideExplanation        *bool    `json:"hide_explanation"`
	TimeLimitMinutes       *int     `json:"time_limit_minutes"`
	IsActive               *bool    `json:"is_active"`
}

type AddCardRequest struct {
	Cards string `json:"cards" validate:"required"`
}

type ImportPDFRequest struct {
	Questions string `json:"questions" validate:"required"`
}
