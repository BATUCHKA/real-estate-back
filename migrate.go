package main

import (
	"fmt"
	"log"
	"reflect"

	"gitlab.com/steppelink/odin/odin-backend/database"
	"gitlab.com/steppelink/odin/odin-backend/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db := database.Database
	db.GormDB.Logger = logger.Default.LogMode(logger.Info) // Enable Logging

	modelsList := []interface{}{
	}

	for _, model := range modelsList {
		err := autoMigrateModel(db.GormDB, model)
		if err != nil {
			log.Fatalf("Failed to auto migrate model %v: %v", reflect.TypeOf(model), err)
		}
	}
}

func autoMigrateModel(db *gorm.DB, model interface{}) error {
	log.Printf("Migratinwg model: %v", reflect.TypeOf(model))
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("auto migration failed for model %v: %w", reflect.TypeOf(model), err)
	}
	log.Printf("Successfully migrated model: %v", reflect.TypeOf(model))
	return nil
}
