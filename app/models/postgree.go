package models

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"not null"`
	Email     string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
	Token     *string        `gorm:"default:NULL"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Task struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	UserId         uint           `gorm:"not null"` 
	Task           string         `gorm:"not null"`
	TaskDate       time.Time      `gorm:"not null"`
	AttachmentFile *string        `gorm:"default:NULL"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// Fungsi untuk inisialisasi database
func InitDb(config map[string]interface{}) (*gorm.DB, error) {
	dbHost, ok := config["dbHost"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for dbHost")
	}
	dbUser, ok := config["dbUser"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for dbUser")
	}
	dbPassword, ok := config["dbPassword"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for dbPassword")
	}
	dbName, ok := config["dbName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for dbName")
	}
	dbPort, ok := config["dbPort"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for dbPort")
	}

	// Membuat DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Koneksi ke database PostgreSQL dengan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Migrasi Model
	// db.AutoMigrate(&Users{},&Task{})

	return db, nil
}