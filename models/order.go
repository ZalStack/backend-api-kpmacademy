package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID                    uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID                uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	User                  User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PackageID             *uuid.UUID     `gorm:"type:char(36)" json:"package_id"`
	Package               *Package       `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	VideoOrderID          *uuid.UUID     `gorm:"type:char(36)" json:"video_order_id"`
	Type                  string         `gorm:"type:varchar(50);default:package;not null" json:"type"`
	OrderNumber           string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	TotalPrice            float64        `gorm:"type:decimal(10,2);default:0" json:"total_price"`
	IsCustomAmount        bool           `gorm:"default:false" json:"is_custom_amount"`
	PaymentStatus         string         `gorm:"type:varchar(50);default:pending;not null" json:"payment_status"`
	TransactionID         string         `gorm:"type:varchar(255);uniqueIndex" json:"transaction_id"`
	PaymentReference      string         `gorm:"type:varchar(255)" json:"payment_reference"`
	PaymentURL            string         `gorm:"type:text" json:"payment_url"`
	PaymentType           string         `gorm:"type:varchar(50)" json:"payment_type"`
	PaymentTime           *time.Time     `json:"payment_time"`
	PaymentNotes          string         `gorm:"type:text" json:"payment_notes"`
	Enrollment            string         `gorm:"type:json" json:"enrollment"`
	MembershipDurationDays *int          `json:"membership_duration_days"`
	MembershipStart       *time.Time     `json:"membership_start"`
	MembershipEnd         *time.Time     `json:"membership_end"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type CreateOrderRequest struct {
	PackageID     string  `json:"package_id"`
	TotalPrice    float64 `json:"total_price"`
	IsCustomAmount bool   `json:"is_custom_amount"`
}

type PaymentCallback struct {
	OrderID       string  `json:"order_id"`
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	PaymentType   string  `json:"payment_type"`
	GatewayData   string  `json:"gateway_data"`
}
