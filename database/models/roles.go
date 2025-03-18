package models

import (
	"database/sql"
	"time"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type RoleKeyType string

const (
	RoleClient     RoleKeyType = "client"
	RoleAgent      RoleKeyType = "agent"
	RoleAccountant RoleKeyType = "accountant"
	RoleAdmin      RoleKeyType = "admin"
)

type Roles struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key       RoleKeyType  `gorm:"unique;not null" json:"key"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}

func RoleFlush() {
	db := database.Database
	types := []RoleKeyType{RoleAdmin, RoleAgent, RoleAccountant, RoleClient}
	for _, v := range types {
		staticTable := &Roles{
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
