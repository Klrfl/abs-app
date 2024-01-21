package handlers

import (
	"abs-app/database"
	"abs-app/models"

	"github.com/gofiber/fiber/v2"
)

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
