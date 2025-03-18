package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Project struct {
	ID                 uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	TotalApartments    int             `json:"total_apartments"`
	RemaningApartments int             `json:"remaining_apartments"`
	SoldApartments     int             `json:"sold_apartments"`
	AdvantagesHTML     string          `json:"advantages_html" gorm:"type:text"`
	Longitude          *float64        `json:"longitude"`
	Latitude           *float64        `json:"latitude"`
	ProjectStartedAt   datatypes.Date  `json:"project_started_at"`
	ProjectEndedAt     *datatypes.Date `json:"project_ended_at"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          sql.NullTime    `gorm:"index" json:"-"`
}
