// available, reserved, sold
package models

import (
	"database/sql"
	"time"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type ApartmentStatusKeyType string

const (
	ApartmentStatusAvailable ApartmentStatusKeyType = "available"
	ApartmentStatusReserved  ApartmentStatusKeyType = "reserved"
	ApartmentStatusSold      ApartmentStatusKeyType = "sold"
)

type ApartmentStatus struct {
	ID        uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key       ApartmentStatusKeyType `gorm:"unique;not null" json:"key"`
	CreatedAt time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt sql.NullTime           `gorm:"index" json:"-"`
}

func ApartmentStatusFlush() {
	db := database.Database
	types := []ApartmentStatusKeyType{ApartmentStatusReserved, ApartmentStatusSold, ApartmentStatusAvailable}
	for _, v := range types {
		staticTable := &ApartmentStatus{
			Key: v,
		}
		db.GormDB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"key": v,
			}),
		}).Create(&staticTable)
	}
}
