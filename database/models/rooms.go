package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID                 uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Area               uint         `json:"area"`
	CreatedAt          time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          sql.NullTime `gorm:"index" json:"-"`
}
