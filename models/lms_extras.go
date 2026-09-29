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

type LeaderboardResponseItem struct {
	Rank          int     `json:"rank"`
	UserID        string  `json:"user_id"`
	Name          string  `json:"name"`
	StudentClass  string  `json:"student_class"`
	SchoolName    string  `json:"school_name"`
	Avatar        string  `json:"avatar"`
	AvgScore      float64 `json:"avg_score"`
	TotalPractice int     `json:"total_practice"`
	TotalCorrect  int     `json:"total_correct"`
}

type Tutor struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Role            string         `gorm:"type:varchar(255);not null" json:"role"`
	Experience      string         `gorm:"type:varchar(100)" json:"experience"`
	Rating          float64        `gorm:"type:decimal(3,2);default:5.0" json:"rating"`
	TotalSessions   int            `gorm:"default:0" json:"total_sessions"`
	Avatar          string         `gorm:"type:text" json:"avatar"`
	PricePerSession float64        `gorm:"type:decimal(12,2);not null" json:"price_per_session"`
	Specialties     string         `gorm:"type:text" json:"specialties_raw,omitempty"`
	AvailableSlots  string         `gorm:"type:text" json:"available_slots_raw,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	SortOrder       int            `gorm:"default:0" json:"sort_order"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tutor) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TutorBooking struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	TutorID           uuid.UUID      `gorm:"type:char(36);index;not null" json:"tutor_id"`
	TutorName         string         `gorm:"type:varchar(255)" json:"tutor_name"`
	StudentID         uuid.UUID      `gorm:"type:char(36);index;not null" json:"student_id"`
	StudentName       string         `gorm:"type:varchar(255);not null" json:"student_name"`
	ScheduledDate     string         `gorm:"type:varchar(20);not null" json:"scheduled_date"`
	ScheduledTimeSlot string         `gorm:"type:varchar(50);not null" json:"scheduled_time_slot"`
	Subject           string         `gorm:"type:varchar(100);not null" json:"subject"`
	Status            string         `gorm:"type:varchar(30);default:'confirmed'" json:"status"`
	MeetingURL        string         `gorm:"type:text" json:"meeting_url"`
	Notes             string         `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (tb *TutorBooking) BeforeCreate(tx *gorm.DB) error {
	if tb.ID == uuid.Nil {
		tb.ID = uuid.New()
	}
	return nil
}

type BookTutorRequest struct {
	ScheduledDate     string `json:"scheduled_date" validate:"required"`
	ScheduledTimeSlot string `json:"scheduled_time_slot" validate:"required"`
	Subject           string `json:"subject" validate:"required"`
	Notes             string `json:"notes"`
}
