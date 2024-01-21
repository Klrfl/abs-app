package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

func GetAllDrinks(c *fiber.Ctx) error {
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

func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order

	result := database.DB.Table("orders").Joins("join members on orders.member_id=members.id join drinks on orders.drink_id=drinks.id").Select("orders.id, members.member_name, drinks.drink_name, drinks.drink_type, drinks.hot_price, drinks.cold_price, orders.created_at").Find(&orders)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
			"data":    orders,
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no orders yet",
			"data":    orders,
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": orders,
	})
}
