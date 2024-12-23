package database

import (
	"fmt"
	"time"

	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func ConnectDatabase() (*Database, error) {
	dbConfig := configs.LoadDatabaseConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=false", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("unable to load env file: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{DB: db}, nil
}

func (db *Database) Migrate() {
	db.DB.AutoMigrate(&models.User{}, &models.Product{})
	fmt.Println("Database migrated")
}
