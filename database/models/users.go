package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Email       string       `json:"email" gorm:"uniqueIndex"`
	Password    string       `json:"password,omitempty"`
	PhoneNumber string       `json:"phone_number"`
	RoleID      string       `json:"role_id"`
	Role        Roles        `json:"-"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   sql.NullTime `gorm:"index" json:"-"`
	// Verified    bool         `gorm:"default:false" json:"-"`
	// EmailVerifySecret           string          `json:"-"`
	// EmailVerifySecretExpireAt   sql.NullTime    `json:"-"`
}
