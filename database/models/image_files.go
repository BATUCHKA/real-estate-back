package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ImageFile struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MimeType    string       `json:"mime_type"`
	FileName    string       `gorm:"index" json:"file_name"`
	FileSize    uint         `json:"file_size"`
	ApartmentID *string      `gorm:"index" json:"apartment_id"`
	Apartment   Apartment    `json:"-"`
	ProjectID   *string      `json:"project_id" gorm:"index"`
	Project     Project      `json:"-"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   sql.NullTime `gorm:"index" json:"-"`
}
