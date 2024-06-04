package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetMenuTypes(c *fiber.Ctx) error {
	var menuTypes []*models.MenuType
	result := database.DB.Find(&menuTypes)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something went wrong when querying database",
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": menuTypes,
	})
}

func GetMenuTypeByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid ID",
		})
	}

	var menuType *models.MenuType
	result := database.DB.
		Where("id = ?", id).
		Find(&menuType).
		Limit(1)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something went wrong when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no menu type found",
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": menuType,
	})
}

func AddNewMenuType(c *fiber.Ctx) error {
	var menuType models.MenuType

	if err := c.BodyParser(&menuType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with request payload",
		})
	}

	result := database.DB.Create(&menuType)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when creating new menu type",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "successfully create new menu type",
	})
}

func DeleteMenuTypeByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid ID",
		})
	}

	tx := database.DB.Begin()
	result := tx.Where("id = ?", id).Delete(&models.MenuType{})

	if result.Error != nil || result.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when deleting menu type",
		})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"err":     false,
		"message": fmt.Sprintf("successfully deleted menu type of ID %d", id),
	})
}

func UpdateMenuType(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid ID",
		})
	}

	var existingMenuType models.MenuType
	var newMenuType models.MenuType

	result := database.DB.
		Where("id = ?", id).
		Find(&existingMenuType).
		Limit(1)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when updating menu type",
		})
	}

	if err := c.BodyParser(&newMenuType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with payload",
		})
	}

	result = database.DB.
		Where("id = ?", id).
		Updates(&newMenuType)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "something wrong when updating menu type",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "successfully updated menu type",
	})
}
