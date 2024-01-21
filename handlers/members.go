package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllMembers(c *fiber.Ctx) error {
	var members []models.Member

	result := database.DB.Find(&members)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no data",
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": members,
	})
}
