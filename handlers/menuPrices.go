package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func InsertNewPrices(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "ID not valid",
		})
	}

	var menuItem models.Menu
	result := database.DB.
		Limit(1).
		Where("id = ?", id).
		Find(&menuItem)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": fmt.Sprintf("menuItem with id %s not found", id),
		})
	}

	newVariantValues := new([]*models.VariantValue)

	if err := c.BodyParser(&newVariantValues); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with your prices data",
		})
	}

	for _, newVariantValue := range *newVariantValues {
		newVariantValue.MenuID = id
	}
	result = database.DB.Create(&newVariantValues)

	menuItem.UpdatedAt = time.Now()
	result2 := database.DB.Updates(&menuItem)

	if result.Error != nil || result2.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new prices",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": fmt.Sprintf("prices for menu item with ID %s sucessfully updated", id),
	})
}

func UpdatePrices(c *fiber.Ctx) error {
	c.Accepts("application/json")
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": fmt.Sprintf("menu id (%s) not valid", id),
		})
	}

	menuItem := new(models.Menu)
	result := database.DB.
		Limit(1).
		Where("id = ?", id).
		Find(&menuItem)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": fmt.Sprintf("menu item with id %s doesn't exist", id),
		})
	}

	incomingVariantValues := new([]*models.InputVariantValue)

	if err := c.BodyParser(&incomingVariantValues); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the body of your request",
		})
	}

	for _, variantValue := range *incomingVariantValues {
		targetVariantValue := new(models.VariantValue)
		result = database.DB.
			Where("menu_id = ? AND option_id = ? AND option_value_id = ?", id, variantValue.OptionID, variantValue.OptionValueID).
			Limit(1).
			Find(&targetVariantValue)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database for prices",
			})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"err":     true,
				"message": fmt.Sprintf("one of prices for menu %s not found. Make sure all prices supplied exists", id),
			})
		}

		var newVariantValue models.VariantValue
		newVariantValue.MenuID = id
		newVariantValue.OptionValueID = variantValue.NewOptionValueID
		newVariantValue.OptionID = variantValue.NewOptionID
		newVariantValue.Price = variantValue.Price

		result = database.DB.
			Model(&models.VariantValue{}).
			Where("menu_id = ? AND option_id = ? AND option_value_id = ?", id, variantValue.OptionID, variantValue.OptionValueID).
			Updates(&newVariantValue)

		menuItem.UpdatedAt = time.Now()
		result2 := database.DB.Updates(&menuItem)

		if result.Error != nil || result2.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err":     true,
				"message": "something went wrong when updating prices",
			})
		}

		if result.RowsAffected == 0 {
			return c.JSON(fiber.Map{
				"err":     true,
				"message": "your update doesn't work lol?",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "prices for menu successfully updated",
	})
}

func DeletePrice(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "ID not valid",
		})
	}

	menuItem := new(models.Menu)
	result := database.DB.
		Limit(1).
		Where("id = ?", id).
		Find(&menuItem)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": fmt.Sprintf("menu item with id %s doesn't exist", id),
		})
	}

	OptionID := c.QueryInt("option_id")
	OptionValueID := c.QueryInt("option_value_id")
	result = database.DB.
		Where("option_id = ? AND option_value_id = ?", OptionID, OptionValueID).
		Delete(&models.VariantValue{}, id)

	menuItem.UpdatedAt = time.Now()
	result2 := database.DB.Updates(&menuItem)

	if result.Error != nil || result2.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when deleting prices of menu item from database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": fmt.Sprintf("failed to delete menu item with id %s", id),
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": fmt.Sprintf("prices of menu item with id %s successfully deleted", id),
	})
}
