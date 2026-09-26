package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Banner struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Tag       string         `gorm:"type:varchar(100)" json:"tag"`
	SubTag    string         `gorm:"type:varchar(100)" json:"sub_tag"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Subtitle  string         `gorm:"type:text" json:"subtitle"`
	ImageURL  string         `gorm:"type:text;not null" json:"image_url"`
	ActionURL string         `gorm:"type:text" json:"action_url"`
	Route     string         `gorm:"type:varchar(255)" json:"route"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Banner) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type LiveClass struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title           string         `gorm:"type:varchar(255);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	TutorName       string         `gorm:"type:varchar(255);not null" json:"tutor_name"`
	TutorRole       string         `gorm:"type:varchar(255)" json:"tutor_role"`
	TutorAvatar     string         `gorm:"type:text" json:"tutor_avatar"`
	Subject         string         `gorm:"type:varchar(100);not null" json:"subject"`
	Grade           string         `gorm:"type:varchar(100);not null" json:"grade"`
	ScheduledAt     time.Time      `json:"scheduled_at"`
	DurationMinutes int            `gorm:"default:60" json:"duration_minutes"`
	ZoomURL         string         `gorm:"type:text;not null" json:"zoom_url"`
	IsLive          bool           `gorm:"default:false" json:"is_live"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (lc *LiveClass) BeforeCreate(tx *gorm.DB) error {
	if lc.ID == uuid.Nil {
		lc.ID = uuid.New()
	}
	return nil
}

type Module struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Subject     string         `gorm:"type:varchar(100);not null" json:"subject"`
	Grade       string         `gorm:"type:varchar(100);not null" json:"grade"`
	Chapter     int            `gorm:"default:1" json:"chapter"`
	TotalPages  int            `gorm:"default:10" json:"total_pages"`
	PDFURL      string         `gorm:"type:text;not null" json:"pdf_url"`
	CoverImage  string         `gorm:"type:text" json:"cover_image"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Module) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type Voucher struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code            string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Title           string         `gorm:"type:varchar(255);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	DiscountPercent int            `gorm:"default:0" json:"discount_percent"`
	DiscountAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	MinPurchase     float64        `gorm:"type:decimal(10,2);default:0" json:"min_purchase"`
	MaxDiscount     float64        `gorm:"type:decimal(10,2);default:0" json:"max_discount"`
	ValidUntil      time.Time      `json:"valid_until"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	UsageLimit      int            `gorm:"default:1000" json:"usage_limit"`
	UsageCount      int            `gorm:"default:0" json:"usage_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Voucher) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

type ClaimVoucherRequest struct {
	Code string `json:"code" validate:"required"`
}
