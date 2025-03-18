package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Key       uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"key"`
	Hash      string       `gorm:"index" json:"hash"`
	Data      string       `json:"data"`
	ExpireAt  time.Time    `gorm:"index" json:"expire_at"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}
