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

	var menu []models.Menu

	var result *gorm.DB
	/* 	joinQuery := "join menu_types on menu.type_id=menu_types.id"
	   	selectQuery := "menu.id, menu.name, menu_types.type, menu.created_at, menu.updated_at" */

	if queries["item_name"] != "" {
		result = database.DB.Where("item_name ILIKE ?", fmt.Sprintf("%%%s%%", queries["item_name"])).Find(&menu)
	} else {
		crosstabJoinQuery := "join crosstab('select menu_id, menu_option_value_id, price from variant_values group by 1,2,3 order by 1,2', $$select distinct menu_option_value_id from variant_values order by 1$$) as ct(menu_id uuid, iced numeric, hot numeric, blend numeric, regular numeric, plain numeric) on menu.id=ct.menu_id"
		menuTypesJoinQuery := "join menu_types on menu.type_id=menu_types.id"
		selectQueryString = "menu.id, menu.name, menu_types.type, ct.iced, ct.hot, ct.blend, ct.regular, ct.plain"
		result = database.DB.Table("menu").Select(selectQueryString).Joins(crosstabJoinQuery).Joins(menuTypesJoinQuery).Where(queries).Find(&menu)
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

	return c.JSON(&menu)
}

func GetMenuItemByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing drink id",
		})
	}

	var menuItem models.BaseMenu

	result := database.DB.First(&menuItem, "id = ?", id)

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

	menuItem := new(models.BaseMenu)

	if err := c.BodyParser(menuItem); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the drink data",
		})
	}

	menuItem.ID = id
	menuItem.Updated_at = time.Now()
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

	result := database.DB.Where("id = ?", id).Delete(&models.BaseMenu{})

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
