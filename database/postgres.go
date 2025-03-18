package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Postgres struct {
	GormDB *gorm.DB
}

var Database *Postgres

func init() {
	log.Println("initializing db connection")
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Print(connectionString)

	sqlDB, err := sql.Open("pgx", connectionString)

	if err != nil {
		log.Panic(err)
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	Database = &Postgres{
		GormDB: gormDB,
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Panic(err)
	}
}
