package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"gorm.io/gorm"
)

func main() {
	db := database.Database
	// db.GormDB.Logger = logger.Default.LogMode(logger.Info)

	// namingStrategy := schema.NamingStrategy{
	// 	SingularTable: true,
	// }

	// db.GormDB.NamingStrategy = namingStrategy

	modelsList := []interface{}{
		&models.User{},
		&models.Session{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Project{},
		&models.Building{},
		&models.Apartment{},
		&models.Room{},
		&models.Config{},
		&models.ImageFile{},
	}

	for _, model := range modelsList {
		err := autoMigrateModel(db.GormDB, model)
		if err != nil {
			log.Fatalf("Failed to auto migrate model %v: %v", reflect.TypeOf(model), err)
		}
	}

	log.Println("All migrations completed successfully with singular table names")
}

func autoMigrateModel(db *gorm.DB, model interface{}) error {
	log.Printf("Migrating model: %v", reflect.TypeOf(model))
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("auto migration failed for model %v: %w", reflect.TypeOf(model), err)
	}
	log.Printf("Successfully migrated model: %v", reflect.TypeOf(model))
	return nil
}
