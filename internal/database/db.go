package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	fmt.Println("Setting up new database connection")

	dbUserName := "postgres"
	dbPassword := "postgres"
	dbHost := "localhost"
	dbName := "postgres"
	dbPort := "5432"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUserName, dbName, dbPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	// if ping fails, we will get an error

	postgresDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := postgresDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// docker run --name product-api-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14-alpine
