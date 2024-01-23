package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var joinQueryString string = "join members on orders.member_id=members.id join drinks on orders.drink_id=drinks.id"
var selectQueryString string = "orders.id, members.member_name, drinks.drink_name, drinks.drink_type, drinks.hot_price, drinks.cold_price, orders.is_completed, orders.created_at"

func GetOrders(c *fiber.Ctx) error {
	var id uuid.UUID
	var err error

	var orders []models.Order

	var result *gorm.DB

	if c.Query("member_id") != "" {
		id, err = uuid.Parse(c.Query("member_id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"err":     true,
				"message": "error when processing member id`",
			})
		}

		result = database.DB.Table("orders").Joins(joinQueryString).Select(selectQueryString).Where("members.id = ?", id).Find(&orders)
	} else {
		result = database.DB.Table("orders").Joins(joinQueryString).Select(selectQueryString).Find(&orders)
	}

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

func GetOrderByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when processing order id",
		})
	}

	order := new(models.Order)
	result := database.DB.Joins(joinQueryString).Select(selectQueryString).Where("orders.id = ?", id).Find(&order)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"message": "no data",
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": order,
	})
}

func CreateNewOrder(c *fiber.Ctx) error {
	newOrder := new(models.BaseOrder)

	if err := c.BodyParser(newOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with your request body",
		})
	}

	result := database.DB.Table("orders").Save(&newOrder)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "new order successfully created",
	})
}

func CompleteOrder(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing id",
		})
	}

	order := new(models.BaseOrder)

	order.Id = id
	order.Completed_at = time.Now()
	result := database.DB.Model(&order).Updates(&order)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "order completed",
	})
}
