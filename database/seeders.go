package database

import (
	"abs-app/models"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Seeder() error {
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

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	newPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 14)
	if err != nil {
		return err
	}

	admins := models.User{
		ID:       uuid.New(),
		Name:     adminUsername,
		Email:    adminEmail,
		Password: string(newPassword), // TODO: replace with secure password from .env
		RoleID:   2,
	}

	DB.Where("role_id = ?", 2).Delete(&models.User{})
	result := DB.Save(&roles)
	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}

	result = DB.Where("email = ?", adminEmail).Save(&admins)
	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
