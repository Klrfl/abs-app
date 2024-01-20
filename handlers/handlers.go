package handlers

import (
	"abs-app/database"
	"abs-app/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllDrinks(c *fiber.Ctx) error {
	queries := c.Queries()

	db := database.GetDBInstance()
	var drinks []models.Drink

	var result *gorm.DB

	if queries["drink_name"] != "" {
		result = db.Where("drink_name ILIKE ?", fmt.Sprintf("%%%s%%", queries["drink_name"])).Find(&drinks)
	} else {
		result = db.Where(queries).Find(&drinks)
	}

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

	var drink models.Drink

	result := database.DB.Find(&drink, "id = ?", id)

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
