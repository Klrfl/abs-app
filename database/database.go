package database

import (
	"abs-app/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

func Init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbname, dbPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("error when connecting to database")
	} else {
		log.Println("Successfully connected to database")
	}

	err = DB.AutoMigrate(
		&models.Menu{},
		&models.Order{},
		&models.User{},
		&models.Role{},
		&models.MenuType{},
		&models.VariantValue{},
		&models.MenuAvailableOption{},
		&models.MenuOptionValue{},
		&models.BaseOrderDetail{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database")
	} else {
		log.Println("Database successfully migrated")
	}

	if err := Seeder(); err != nil {
		log.Fatal("failed to seed data")
	}
	log.Println("database seeded")
}
