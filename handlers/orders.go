package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOrders(c *fiber.Ctx) error {
	var id uuid.UUID
	var err error

	var orders []*models.Order

	var result *gorm.DB

	if c.Query("member_id") != "" {
		id, err = uuid.Parse(c.Query("member_id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"err":     true,
				"message": "error when processing member id`",
			})
		}

		result = database.DB.
			Preload("Member").
			Where("members.id = ?", id).
			Find(&orders)
	} else {
		result = database.DB.
			Preload("Member").
			Find(&orders)
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

	// todo: make raw sql query to get []orderDetails

	for _, order := range orders {
		var orderDetails []*models.OrderDetail
		rows, err := database.DB.
			Model(&models.OrderDetail{}).
			Select("order_details.order_id, menu.id, menu.name, menu_types.type, menu_available_options.id, menu_available_options.option, menu_option_values.id, menu_option_values.value, order_details.quantity, order_details.quantity * variant_values.price as total_price").
			Joins("join variant_values on order_details.menu_id=variant_values.menu_id and order_details.menu_option_value_id=variant_values.option_value_id").
			Joins("join menu on order_details.menu_id=menu.id").
			Joins("join menu_types on menu.type_id=menu_types.id").
			Joins("join menu_option_values on order_details.menu_option_value_id=menu_option_values.id").
			Joins("join menu_available_options on order_details.menu_option_id=menu_available_options.id").
			Where("order_details.order_id = ?", order.ID).
			Rows()

		defer rows.Close()

		if err != nil {
			return c.JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database for order_details",
			})
		}

		for rows.Next() {
			var orderDetail models.OrderDetail
			rows.Scan(
				&orderDetail.OrderID,
				&orderDetail.MenuID,
				&orderDetail.MenuName,
				&orderDetail.MenuType,
				&orderDetail.MenuOptionID,
				&orderDetail.MenuOption,
				&orderDetail.MenuOptionValueID,
				&orderDetail.MenuOptionValue,
				&orderDetail.Quantity,
				&orderDetail.TotalPrice,
			)
			orderDetails = append(orderDetails, &orderDetail)
		}

		order.OrderDetails = orderDetails
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
	result := database.DB.Preload("Member").Where("orders.id = ?", id).Find(&order)

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

	var orderDetails []*models.OrderDetail
	rows, err := database.DB.
		Model(&models.OrderDetail{}).
		Select("order_details.order_id, menu.id, menu.name, menu_types.type, menu_available_options.id, menu_available_options.option, menu_option_values.id, menu_option_values.value, order_details.quantity, order_details.quantity * variant_values.price as total_price").
		Joins("join variant_values on order_details.menu_id=variant_values.menu_id and order_details.menu_option_value_id=variant_values.option_value_id").
		Joins("join menu on order_details.menu_id=menu.id").
		Joins("join menu_types on menu.type_id=menu_types.id").
		Joins("join menu_option_values on order_details.menu_option_value_id=menu_option_values.id").
		Joins("join menu_available_options on order_details.menu_option_id=menu_available_options.id").
		Where("order_details.order_id = ?", order.ID).
		Rows()

	defer rows.Close()

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database for order_details",
		})
	}

	for rows.Next() {
		var orderDetail models.OrderDetail
		rows.Scan(
			&orderDetail.OrderID,
			&orderDetail.MenuID,
			&orderDetail.MenuName,
			&orderDetail.MenuType,
			&orderDetail.MenuOptionID,
			&orderDetail.MenuOption,
			&orderDetail.MenuOptionValueID,
			&orderDetail.MenuOptionValue,
			&orderDetail.Quantity,
			&orderDetail.TotalPrice,
		)
		orderDetails = append(orderDetails, &orderDetail)
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": order,
	})
}

func CreateNewOrder(c *fiber.Ctx) error {
	newOrder := new(models.Order)

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

	order := new(models.Order)

	order.ID = id
	order.CompletedAt = time.Now()
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
