package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	ID           uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RoleID       string       `json:"role_id"`
	Role         Roles        `json:"-"`
	PermissionID string       `json:"permission_id"`
	Permission   Permission   `json:"-"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    sql.NullTime `gorm:"index" json:"-"`
}
