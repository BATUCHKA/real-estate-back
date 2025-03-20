package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Building struct {
	ID              uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string          `json:"name"`
	Number          string          `json:"number"`
	ProjectID       string          `json:"project_id" gorm:"index"`
	Project         Project         `json:"-"`
	Description     string          `json:"description"`
	TotalFloors     int             `json:"total_floors"`
	TotalApartments int             `json:"total_apartments"`
	CompletionDate  *datatypes.Time `json:"completion_date"`
	StartDate       *datatypes.Time `json:"start_date"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       sql.NullTime    `gorm:"index" json:"-"`
}
