package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Apartment struct {
	ID              uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Area            uint            `json:"area"`// talbai
	Floor           uint            `json:"floor"`
	ApartmentNumber string          `json:"apartment_number"`
	Price           float64         `json:"price"`
	StatusID        string          `json:"status_id" gorm:"index"`
	Status          ApartmentStatus `json:"-"`

	BuildingID string   `json:"building_id" gorm:"index"`
	Building   Building `json:"-"`

	ProjectID string  `json:"project_id" gorm:"index"`
	Project   Project `json:"-"`

	Features datatypes.JSON `json:"features"`

	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}
