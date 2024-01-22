package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateNewDrink(c *fiber.Ctx) error {
	var newDrink models.Drink

	if err := c.BodyParser(&newDrink); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with new drink data",
		})
	}

	result := database.DB.Save(&newDrink)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "new drink successfully created",
	})
}

func GetDrinks(c *fiber.Ctx) error {
	queries := c.Queries()

	var drinks []models.Drink

	var result *gorm.DB

	if queries["drink_name"] != "" {
		result = database.DB.Where("drink_name ILIKE ?", fmt.Sprintf("%%%s%%", queries["drink_name"])).Find(&drinks)
	} else {
		result = database.DB.Where(queries).Find(&drinks)
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
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing drink id",
		})
	}

	var drink models.Drink

	result := database.DB.First(&drink, "id = ?", id)

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

func UpdateDrink(c *fiber.Ctx) error {
	c.Accepts("application/json")
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error parsing ID",
		})
	}

	drink := new(models.Drink)

	if err := c.BodyParser(drink); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the drink data",
		})
	}

	drink.Id = id
	drink.Updated_at = time.Now()
	result := database.DB.Save(&drink)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no drink found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "Drink data successfully updated",
	})
}

func DeleteDrink(c *fiber.Ctx) error {
	id := c.Params("id")

	type Drink models.Drink

	result := database.DB.Where("id = ?", id).Delete(&Drink{})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "drink successfully deleted",
	})
}
