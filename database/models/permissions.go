package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Permission struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key         string       `gorm:"uniqueIndex" json:"key"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   sql.NullTime `gorm:"index" json:"-"`
}
