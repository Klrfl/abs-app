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
	"gorm.io/gorm"
)

func CreateNewMenuItem(c *fiber.Ctx) error {
	newMenuItem, incomingMenuItem := new(models.Menu), new(models.Menu)

	if err := c.BodyParser(&incomingMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing request payload: make sure all fields are correct",
		})
	}

	newMenuItem.ID = uuid.New()
	newMenuItem.Name = incomingMenuItem.Name
	newMenuItem.TypeID = incomingMenuItem.TypeID

	tx := database.DB.Begin()
	error := tx.
		Select("Name", "TypeID").
		Create(&newMenuItem).Error

	if error != nil {
		tx.Rollback()

		if pgError := error.(*pgconn.PgError); errors.Is(error, pgError) {
			if pgError.Code == "23503" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"err":     true,
					"message": "make sure type_id is correct",
				})
			}
		}

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

	result := tx.Create(&newVariantValues)

	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()

		if pgError := result.Error.(*pgconn.PgError); errors.Is(result.Error, pgError) {
			if pgError.Code == "23503" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"err":     true,
					"message": "make sure all values supplied are valid",
				})
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting new menu data to database",
		})
	}

	tx.Commit()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "new menu item successfully created",
	})
}

func paginate(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 0)

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GetMenu(c *fiber.Ctx) error {
	var menu []*models.Menu
	var result *gorm.DB

	queryMenuName := c.Query("name")
	queryMenuTypeID := c.QueryInt("type_id")

	if queryMenuName != "" {
		result = database.DB.
			Preload("Type").
			Scopes(paginate(c)).
			Where("menu.name ILIKE ?", fmt.Sprintf("%%%s%%", queryMenuName)).
			Find(&menu)
	} else {
		result = database.DB.
			Preload("Type").
			Scopes(paginate(c)).
			Where(&models.Menu{TypeID: queryMenuTypeID}).
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
		"err":   false,
		"data":  menu,
		"count": len(menu),
	})
}

func GetMenuItemByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing menu item ID",
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

	return c.JSON(fiber.Map{
		"err":  false,
		"data": &menuItem,
	})
}

func UpdateMenuItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error parsing ID",
		})
	}

	incomingMenuItem, existingMenuItem := new(models.Menu), new(models.Menu)

	if err := c.BodyParser(&incomingMenuItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the drink data",
		})
	}

	result := database.DB.
		Where("id = ?", id).
		First(&existingMenuItem)

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

	tx := database.DB.Begin()

	incomingMenuItem.UpdatedAt = time.Now()
	error1 := tx.
		Model(&models.Menu{}).
		Where("id = ?", id).
		Updates(&incomingMenuItem).Error

	newVariantValues := incomingMenuItem.VariantValues
	for _, newVariantValue := range newVariantValues {
		newVariantValue.MenuID = existingMenuItem.ID
	}

	error2 := tx.
		Model(&models.VariantValue{}).
		Where("menu_id = ?", id).
		Updates(&newVariantValues).Error

	if error1 != nil || error2 != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when updating menu item",
		})
	}

	tx.Commit()
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

	tx := database.DB.Begin()

	result := tx.
		Select("VariantValues").
		Delete(&models.Menu{ID: id})

	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "delete didn't work - item already deleted or does not exist",
		})
	}

	tx.Commit()
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

	tx := database.DB.Begin()

	// delete corresponding variant_values first
	// there must be a better way to do this
	err := tx.Delete(&models.VariantValue{}, menuIDs).Error
	if err != nil {
		tx.Rollback()

		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when deleting menu items from database",
		})
	}

	result := tx.Delete(&models.Menu{}, menuIDs)

	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "something went wrong when deleting menu items from database",
		})
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"err":     false,
		"message": "menu items successfully deleted",
	})
}
