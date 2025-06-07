package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"backend/internal/config"
	"backend/internal/models"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Connect to SQLite database
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Contest{},
		&models.Submission{},
		&models.SubmissionResult{},
	)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// Create a new database transaction
func BeginTx() (*gorm.DB, error) {
	return DB.Begin()
}

// Commit a transaction
func CommitTx(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Rollback a transaction
func RollbackTx(tx *gorm.DB) error {
	return tx.Rollback().Error
} 