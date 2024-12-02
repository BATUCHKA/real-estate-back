package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AdminUser struct {
	ID           uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email        string       `json:"email"`
	Verified     bool         `gorm:"default:false" json:"-"`
	IsSuperAdmin bool         `gorm:"default:false" json:"-"`
	Password     string       `json:"password"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    sql.NullTime `gorm:"index" json:"-"`
	// FirstName    string       `json:"first_name"`
	// LastName     string       `json:"last_name"`
	// RoleID         string       `json:"role_id"`
	// Role           Role         `json:"-"`
	// OTPCode        string       `json:"otp_code"`
	// OTPExpiry      time.Time    `json:"otp_expiry"`
	// OTPVerifiedAt  time.Time    `json:"otp_verified_at"`
	// PhoneNumber    string       `gorm:"index" json:"phone_number"`
	// Position       string       `json:"position"`
	// ProfileImageID *string      `gorm:"index" json:"profile_image_id"`
	// ProfileImage   Image        `json:"-"`
}
