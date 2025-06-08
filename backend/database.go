package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	var DB *gorm.DB
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		err = fmt.Errorf("failed to load .env: %w", err)
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("failed to open postgres DB: %w", err)
		return nil, err
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		err = fmt.Errorf("failed to migrate DB: %w", err)
		return nil, err
	}

	return DB, nil
}
