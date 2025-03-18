package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Blocks struct {
	ID                 uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	TotalApartments    int          `json:"total_apartments"`
	RemaningApartments int          `json:"remaining_apartments"`
	SoldApartments     int          `json:"sold_apartments"`
	CreatedAt          time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          sql.NullTime `gorm:"index" json:"-"`
}
