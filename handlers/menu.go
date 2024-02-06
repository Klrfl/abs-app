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

func CreateNewMenuItem(c *fiber.Ctx) error {
	var incomingMenuItem models.InputMenu

	if err := c.BodyParser(&incomingMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with new drink data",
		})
	}

	var newMenuItem models.Menu
	newMenuItem.ID = uuid.New()
	newMenuItem.Name = incomingMenuItem.Name
	newMenuItem.TypeID = incomingMenuItem.TypeID

	result := database.DB.Save(&newMenuItem)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new menu data to database",
		})
	}

	var newVariantValue models.VariantValue
	newVariantValue.MenuID = newMenuItem.ID
	newVariantValue.OptionID = incomingMenuItem.OptionID
	newVariantValue.OptionValueID = incomingMenuItem.OptionValueID
	newVariantValue.Price = incomingMenuItem.Price

	result = database.DB.Save(&newVariantValue)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new menu data to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "new menu item successfully created",
	})
}

func GetMenu(c *fiber.Ctx) error {
	queries := c.Queries()

	var menu []*models.Menu
	var result *gorm.DB

	if queries["name"] != "" {
		result = database.DB.
			Preload("Type").
			Where("menu.name ILIKE ?", fmt.Sprintf("%%%s%%", queries["name"])).
			Find(&menu)
	} else {
		result = database.DB.
			Preload("Type").
			Where(&models.Menu{TypeID: c.QueryInt("type_id")}).
			Find(&menu)
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

	// append price data to each menu item

	for _, menuItem := range menu {
		var variantValues []*models.VariantValue

		err := database.DB.
			Model(&models.VariantValue{}).
			Preload("Option").
			Preload("OptionValue").
			Where("variant_values.menu_id = ?", menuItem.ID).
			Find(&variantValues).Error

		if err != nil {
			return c.JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database",
			})
		}

		menuItem.VariantValues = variantValues
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": menu,
	})
}

func GetMenuItemByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing drink id",
		})
	}

	var menuItem *models.Menu
	result := database.DB.
		Preload("Type").
		Where("id = ?", id).
		Limit(1).
		Find(&menuItem)

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

	var variantValues []*models.VariantValue
	err = database.DB.
		Model(&models.VariantValue{}).
		Preload("Option").
		Preload("OptionValue").
		Where("variant_values.menu_id = ?", menuItem.ID).
		Find(&variantValues).Error

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	menuItem.VariantValues = variantValues

	return c.JSON(&menuItem)
}

func UpdateMenuItem(c *fiber.Ctx) error {
	c.Accepts("application/json")
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error parsing ID",
		})
	}

	menuItem := new(models.Menu)

	if err := c.BodyParser(&menuItem); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the drink data",
		})
	}

	menuItem.ID = id
	menuItem.UpdatedAt = time.Now()
	result := database.DB.Save(&menuItem)

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

func DeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.DB.Where("id = ?", id).Delete(&models.Menu{})

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
