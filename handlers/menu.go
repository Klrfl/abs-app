package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateNewMenuItem(c *fiber.Ctx) error {
	var incomingMenuItem models.Menu

	if err := c.BodyParser(&incomingMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the new menu item's data",
		})
	}

	var newMenuItem models.Menu

	newMenuItem.ID = uuid.New()
	newMenuItem.Name = incomingMenuItem.Name
	newMenuItem.TypeID = incomingMenuItem.TypeID

	error := database.DB.
		Select("Name", "TypeID").
		Create(&newMenuItem).Error

	if error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new menu data to database",
		})
	}

	newVariantValues := incomingMenuItem.VariantValues
	// assign id first don't forget about it
	for _, newVariantValue := range newVariantValues {
		newVariantValue.MenuID = newMenuItem.ID
	}
	error = database.DB.Create(&newVariantValues).Error

	if error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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

	var incomingMenuItem models.Menu

	if err := c.BodyParser(&incomingMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the drink data",
		})
	}

	var menuItem models.Menu
	result := database.DB.Where("id = ?", id).First(&menuItem)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no menu item found",
		})
	}

	incomingMenuItem.UpdatedAt = time.Now()
	error1 := database.DB.
		Model(&models.Menu{}).
		Where("id = ?", id).
		Updates(&incomingMenuItem).Error

	newVariantValues := incomingMenuItem.VariantValues
	for _, newVariantValue := range newVariantValues {
		newVariantValue.MenuID = menuItem.ID
	}

	error2 := database.DB.
		Model(&models.VariantValue{}).
		Where("menu_id = ?", id).
		Updates(&newVariantValues).Error

	if error1 != nil || error2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when updating menu item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "Drink data successfully updated",
	})
}

func DeleteMenuItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "ID not valid",
		})
	}

	result := database.DB.
		Select("VariantValues").
		Delete(&models.Menu{ID: id})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "delete didn't work - item already deleted or does not exist",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "menu item successfully deleted",
	})
}

func DeleteMenuItems(c *fiber.Ctx) error {
	var menuIDs uuid.UUIDs

	if err := c.BodyParser(&menuIDs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with your request body",
		})
	}

	result := database.DB.Delete(&models.Menu{}, menuIDs)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something went wrong when deleting menu items from database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "delete doesn't work",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "menu items successfully deleted",
	})
}
