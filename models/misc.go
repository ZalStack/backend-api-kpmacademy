package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginLog struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"` 
	Location  string    `gorm:"type:varchar(255)" json:"location"`
	LoginAt   time.Time `json:"login_at"`	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (l *LoginLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type Notification struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      string     `gorm:"type:varchar(100);not null" json:"type"`
	Title     string     `gorm:"type:varchar(255);not null" json:"title"`
	Message   string     `gorm:"type:text;not null" json:"message"`
	Data      string     `gorm:"type:json" json:"data"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type ContactForm struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID         *uuid.UUID     `gorm:"type:char(36)" json:"user_id"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Email          string         `gorm:"type:varchar(255);not null" json:"email"`
	Phone          string         `gorm:"type:varchar(20)" json:"phone"`
	Subject        string         `gorm:"type:varchar(255)" json:"subject"`
	Message        string         `gorm:"type:text;not null" json:"message"`
	RecipientEmail string         `gorm:"type:varchar(255);default:sekretariat@bpi.or.id" json:"recipient_email"`
	Status         string         `gorm:"type:enum('pending','read','replied','closed');default:pending" json:"status"`
	AdminReply     string         `gorm:"type:text" json:"admin_reply"`
	RepliedAt      *time.Time     `json:"replied_at"`
	ReadAt         *time.Time     `json:"read_at"`
	IPAddress      string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent      string         `gorm:"type:text" json:"user_agent"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (cf *ContactForm) BeforeCreate(tx *gorm.DB) error {
	if cf.ID == uuid.Nil {
		cf.ID = uuid.New()
	}
	return nil
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID string    `gorm:"type:varchar(100)" json:"session_id"`
	UserID    *uuid.UUID `gorm:"type:char(36)" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string    `gorm:"type:enum('user','assistant');not null" json:"role"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	IsAI      bool      `gorm:"default:false" json:"is_ai"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	return nil
}

type PasswordResetToken struct {
	Email     string    `gorm:"type:varchar(255);primaryKey" json:"email"`
	TokenHash string    `gorm:"type:varchar(64);not null" json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SubmitContactRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message" validate:"required"`
}

type ReplyContactRequest struct {
	AdminReply string `json:"admin_reply" validate:"required"`
}

type UpdateContactStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type SendChatRequest struct {
	Message   string `json:"message" validate:"required"`
	SessionID string `json:"session_id"`
}
