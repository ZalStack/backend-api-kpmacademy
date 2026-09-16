package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title            string         `gorm:"type:varchar(255);not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	PackageID        *uuid.UUID     `gorm:"type:char(36)" json:"package_id"`
	Package          *Package       `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	VideoFile        string         `gorm:"type:varchar(255)" json:"video_file"`
	VideoURL         string         `gorm:"type:varchar(500)" json:"video_url"`
	Thumbnail        string         `gorm:"type:varchar(255)" json:"thumbnail"`
	Price            float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	DiscountType     string         `gorm:"type:varchar(20)" json:"discount_type"`
	DiscountValue    *float64       `gorm:"type:decimal(10,2)" json:"discount_value"`
	AccessDurationDays int          `gorm:"default:30" json:"access_duration_days"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	IsPayWhatYouWant bool           `gorm:"default:false" json:"is_pay_what_you_want"`
	MinPayAmount     float64        `gorm:"type:decimal(10,2);default:0" json:"min_pay_amount"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

type VideoOrder struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	User             User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VideoID          uuid.UUID      `gorm:"type:char(36);not null" json:"video_id"`
	Video            Video          `gorm:"foreignKey:VideoID" json:"video,omitempty"`
	OrderNumber      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	TotalPrice       float64        `gorm:"type:decimal(10,2);default:0" json:"total_price"`
	PaymentStatus    string         `gorm:"type:varchar(50);default:pending;not null" json:"payment_status"`
	TransactionID    string         `gorm:"type:varchar(255);uniqueIndex" json:"transaction_id"`
	PaymentReference string         `gorm:"type:varchar(255)" json:"payment_reference"`
	PaymentURL       string         `gorm:"type:text" json:"payment_url"`
	PaymentType      string         `gorm:"type:varchar(50)" json:"payment_type"`
	PaymentTime      *time.Time     `json:"payment_time"`
	AccessGranted    bool           `gorm:"default:false" json:"access_granted"`
	AccessStart      *time.Time     `json:"access_start"`
	AccessEnd        *time.Time     `json:"access_end"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (vo *VideoOrder) BeforeCreate(tx *gorm.DB) error {
	if vo.ID == uuid.Nil {
		vo.ID = uuid.New()
	}
	return nil
}

type CreateVideoRequest struct {
	Title              string  `json:"title" validate:"required"`
	Description        string  `json:"description"`
	PackageID          string  `json:"package_id"`
	VideoFile          string  `json:"video_file"`
	VideoURL           string  `json:"video_url"`
	Thumbnail          string  `json:"thumbnail"`
	Price              float64 `json:"price"`
	DiscountType       string  `json:"discount_type"`
	DiscountValue      float64 `json:"discount_value"`
	AccessDurationDays int     `json:"access_duration_days"`
	IsPayWhatYouWant   bool    `json:"is_pay_what_you_want"`
	MinPayAmount       float64 `json:"min_pay_amount"`
}

type UpdateVideoRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	PackageID          string   `json:"package_id"`
	VideoFile          string   `json:"video_file"`
	VideoURL           string   `json:"video_url"`
	Thumbnail          string   `json:"thumbnail"`
	Price              *float64 `json:"price"`
	DiscountType       string   `json:"discount_type"`
	DiscountValue      *float64 `json:"discount_value"`
	AccessDurationDays *int     `json:"access_duration_days"`
	IsActive           *bool    `json:"is_active"`
	IsPayWhatYouWant   *bool    `json:"is_pay_what_you_want"`
	MinPayAmount       *float64 `json:"min_pay_amount"`
}

type CreateVideoOrderRequest struct {
	TotalPrice float64 `json:"total_price"`
}
