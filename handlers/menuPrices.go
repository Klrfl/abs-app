package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetItemPrices(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	optionId := c.QueryInt("option_id", 0)
	optionValueId := c.QueryInt("option_value_id", 0)
	// quantity := c.QueryInt("option_value_id", 0)

	if optionId == 0 || optionValueId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "option_id and option_value_id is required",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error parsing menu ID",
		})
	}

	var prices models.VariantValue
	result := database.DB.
		Preload("Option").
		Preload("OptionValue").
		Where("menu_id = ?", id).
		Where("option_id = ?", optionId).
		Where("option_value_id = ?", optionValueId).
		Limit(1).
		Find(&prices)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when finding prices",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     true,
			"message": "prices not found",
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": prices,
	})
}

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

	tx := database.DB.Begin()
	result = tx.Create(&newVariantValues)

	if result.Error != nil {
		tx.Rollback()

		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23505" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"err":     true,
					"message": "one of the items contains a combination of option_id and option_value_id that already exists",
				})
			}

			if pgError.Code == "23503" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"err":     true,
					"message": "make sure both option_id and option_value_id are valid",
				})
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new prices",
		})
	}

	menuItem.UpdatedAt = time.Now()
	result = tx.Updates(&menuItem)

	if result.Error != nil {
		tx.Rollback()
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new prices",
		})
	}

	tx.Commit()
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

		tx := database.DB.Begin()
		result = tx.
			Model(&models.VariantValue{}).
			Where("menu_id = ? AND option_id = ? AND option_value_id = ?", id, variantValue.OptionID, variantValue.OptionValueID).
			Updates(&newVariantValue)

		menuItem.UpdatedAt = time.Now()
		result2 := tx.Updates(&menuItem)

		if result.Error != nil || result2.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err":     true,
				"message": "something went wrong when updating prices",
			})
		}

		tx.Commit()
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
	menuID, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "ID not valid",
		})
	}

	menuPrices := new(models.VariantValue)

	if err := c.BodyParser(&menuPrices); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing request payload",
		})
	}

	menuItem := new(models.Menu)

	result := database.DB.
		Limit(1).
		Where("id = ?", menuID).
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
			"message": fmt.Sprintf("menu item with ID %s doesn't exist", menuID),
		})
	}

	tx := database.DB.Begin()
	result = tx.
		Where("option_id = ? AND option_value_id = ?", menuPrices.OptionID, menuPrices.OptionValueID).
		Delete(&models.VariantValue{}, menuID)

	menuItem.UpdatedAt = time.Now()
	result2 := database.DB.Updates(&menuItem)

	if result.Error != nil || result2.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when deleting prices of menu item from database",
		})
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "combination of option_id and option_value_id already deleted",
		})
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"err":     false,
		"message": fmt.Sprintf("prices of menu item with id %s successfully deleted", menuID),
	})
}
