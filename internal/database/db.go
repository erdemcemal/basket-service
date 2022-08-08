package database

import (
	"fmt"
	log "github.com/siruspen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// NewDatabase - creates a new database connection.
func NewDatabase() (*gorm.DB, error) {
	log.Info("Setting up new database connection")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("Error opening database connection: %s", err)
		return db, fmt.Errorf("failed to connect to database: %w", err)
	}

	postgresDB, err := db.DB()
	if err != nil {
		log.Errorf("Error getting database connection: %s", err)
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		log.Errorf("Error pinging database connection: %s", err)
		return nil, err
	}
	return db, nil
}
