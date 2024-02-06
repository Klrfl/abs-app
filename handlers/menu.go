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

var (
	availableOptionsJoinString string = "join menu_available_options on variant_values.option_id=menu_available_options.id"
	menuOptionValuesJoinString string = "join menu_option_values on variant_values.option_value_id=menu_option_values.id"
	selectQueryString          string = "menu.id, menu.name, menu_types.type, ct.iced, ct.hot, ct.blend, ct.regular, ct.plain, menu.created_at, menu.updated_at"
	variantValuesQueryString   string = "variant_values.menu_id, variant_values.option_id, menu_available_options.option, variant_values.option_value_id, menu_option_values.value as option_value, price"
)

func CreateNewMenuItem(c *fiber.Ctx) error {
	var newMenuItem models.BaseMenu

	if err := c.BodyParser(&newMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with new drink data",
		})
	}

	result := database.DB.Save(&newMenuItem)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
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
		result = database.DB.Preload("Type").Where("name ILIKE ?", fmt.Sprintf("%%%s%%", queries["name"])).Find(&menu)
	} else {
		result = database.DB.Preload("Type").Find(&menu)
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

		rows, err := database.DB.
			Model(&models.VariantValue{}).
			Select(variantValuesQueryString).
			Joins(availableOptionsJoinString).
			Joins(menuOptionValuesJoinString).
			Where("variant_values.menu_id = ?", menuItem.ID).
			Rows()

		defer rows.Close()

		if err != nil {
			return c.JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database",
			})
		}

		for rows.Next() {
			var variantValue models.VariantValue
			rows.Scan(
				&variantValue.MenuID,
				&variantValue.OptionID,
				&variantValue.Option,
				&variantValue.OptionValueID,
				&variantValue.OptionValue,
				&variantValue.Price,
			)
			variantValues = append(variantValues, &variantValue)
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
	result := database.DB.Preload("Type").Where("id = ?", id).First(&menuItem)

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

	rows, err := database.DB.
		Model(&models.VariantValue{}).
		Select(variantValuesQueryString).
		Joins(availableOptionsJoinString).
		Joins(menuOptionValuesJoinString).
		Where("variant_values.menu_id = ?", menuItem.ID).
		Rows()

	defer rows.Close()

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	for rows.Next() {
		var variantValue models.VariantValue
		rows.Scan(
			&variantValue.MenuID,
			&variantValue.OptionID,
			&variantValue.Option,
			&variantValue.OptionValueID,
			&variantValue.OptionValue,
			&variantValue.Price,
		)
		variantValues = append(variantValues, &variantValue)
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

	if err := c.BodyParser(menuItem); err != nil {
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
