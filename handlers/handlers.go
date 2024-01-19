package handlers

import (
	"abs-app/database"
	"abs-app/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllDrinks(c *fiber.Ctx) error {
	db := database.GetDBInstance()
	var drinks []models.Drink

	result := db.Find(&drinks)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "query failed",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no data",
		})
	}

	return c.JSON(&drinks)
}

func GetDrinkByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.GetDBInstance()
	var drink models.Drink

	result := db.Find(&drink, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "query failed",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no data",
		})
	}
	return c.JSON(&drink)
}
