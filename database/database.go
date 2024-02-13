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
		log.Println("Successfully connected to database")
		log.Fatal(err)
	}

	log.Println("Successfully connected to database")

	err = DB.AutoMigrate(
		&models.Menu{},
		&models.Order{},
		&models.User{},
		&models.Role{},
		&models.MenuType{},
		&models.VariantValue{},
		&models.MenuAvailableOption{},
		&models.MenuOptionValue{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database")
	}

	Seeder()
	log.Println("Database successfully migrated")
}
