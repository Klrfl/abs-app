package database

import (
	"abs-app/models"
	"log"

	"github.com/google/uuid"
)

func Seeder() {
	roles := []models.Role{
		{
			ID:   1,
			Name: "user",
		},
		{
			ID:   2,
			Name: "admin",
		},
	}

	admins := []models.User{
		{
			ID:       uuid.New(),
			Name:     "Admin",
			Email:    "admin@admin.com",
			Password: "admin", // TODO: replace with secure password from .env
			RoleID:   2,
		},
	}

	result := DB.Save(&roles)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Fatal("failed to seed data")
	}

	result = DB.Save(&admins)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Fatal("failed to seed data")
	}
}
