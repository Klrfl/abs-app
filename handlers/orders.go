package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getOrderDetails(orderID uuid.UUID) (*sql.Rows, error) {
	rows, err := database.DB.
		Model(&models.OrderDetail{}).
		Select("order_details.order_id, menu.id as menu_id, menu.name as menu_name, menu_types.type as menu_type, menu_available_options.id as menu_option_id, menu_available_options.option as menu_option, menu_option_values.id as menu_option_value_id, menu_option_values.value as menu_option_value, order_details.quantity, order_details.quantity * variant_values.price as total_price").
		Joins("join variant_values on order_details.menu_id=variant_values.menu_id and order_details.menu_option_value_id=variant_values.option_value_id").
		Joins("join menu on order_details.menu_id=menu.id").
		Joins("join menu_types on menu.type_id=menu_types.id").
		Joins("join menu_option_values on order_details.menu_option_value_id=menu_option_values.id").
		Joins("join menu_available_options on order_details.menu_option_id=menu_available_options.id").
		Where("order_details.order_id = ?", orderID).
		Rows()

	return rows, err
}

func GetPendingOrders(c *fiber.Ctx) error {
	var orders []*models.Order

	var result *gorm.DB

	if c.Query("user_id") != "" {
		orderID, err := uuid.Parse(c.Query("user_id"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"err":     true,
				"message": "error when processing user ID",
			})
		}

		result = database.DB.
			Preload("User").
			Preload("User.Role").
			Where("users.id = ? AND is_completed = ?", orderID, false).
			Find(&orders)
	} else {
		result = database.DB.
			Preload("User").
			Preload("User.Role").
			Where("is_completed = ?", false).
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

	for _, order := range orders {
		var orderDetails []*models.OrderDetail

		rows, err := getOrderDetails(order.ID)
		defer rows.Close()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database for order details",
			})
		}

		for rows.Next() {
			var orderDetail models.OrderDetail
			database.DB.ScanRows(rows, &orderDetail)
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
	result := database.DB.
		Preload("User").
		Where("orders.id = ?", id).
		Find(&order)

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

	var orderDetails []*models.OrderDetail
	rows, err := getOrderDetails(order.ID)
	defer rows.Close()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database for order details",
		})
	}

	for rows.Next() {
		var orderDetail models.OrderDetail
		database.DB.ScanRows(rows, &orderDetail)
		orderDetails = append(orderDetails, &orderDetail)
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": order,
	})
}

func GetOrdersForUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	var orders []*models.Order

	result := database.DB.
		Preload("User").
		Preload("User.Role").
		Where("user_id = ? AND is_completed = ?", userID, false).
		Find(&orders)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     true,
			"message": "no orders yet",
		})
	}

	for _, order := range orders {
		rows, err := getOrderDetails(order.ID)
		defer rows.Close()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err":     true,
				"message": "error when querying database for order details",
			})
		}

		var orderDetails []*models.OrderDetail
		for rows.Next() {
			var orderDetail models.OrderDetail
			database.DB.ScanRows(rows, &orderDetail)
			orderDetails = append(orderDetails, &orderDetail)
		}

		order.OrderDetails = orderDetails
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": orders,
	})
}

func GetOrdersForUserByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid order ID",
		})
	}

	userID := c.Locals("user_id")

	order := new(models.Order)
	result := database.DB.
		Limit(1).
		Where("user_id = ? AND id = ?", userID, orderID).
		Preload("User").
		Preload("User.Role").
		Find(&order)

	rows, err := getOrderDetails(order.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when quering database for order details",
		})
	}

	var orderDetails []*models.OrderDetail
	for rows.Next() {
		var orderDetail models.OrderDetail
		database.DB.ScanRows(rows, &orderDetail)
		orderDetails = append(orderDetails, &orderDetail)
	}

	order.OrderDetails = orderDetails

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": fmt.Sprintf("order with ID %s doesn't exist", orderID),
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": order,
	})
}

func CreateNewOrder(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing user ID",
		})
	}

	newOrder, incomingOrder := new(models.Order), new(models.Order)

	if err := c.BodyParser(&incomingOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with your request body",
		})
	}

	newOrder.ID = uuid.New()
	newOrder.UserID = userID

	result := database.DB.Create(&newOrder)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when placing new order",
		})
	}

	var newOrderDetails []models.BaseOrderDetail

	for _, incomingOrderDetail := range incomingOrder.OrderDetails {
		var newOrderDetail models.BaseOrderDetail
		newOrderDetail.OrderID = newOrder.ID
		newOrderDetail.MenuID = incomingOrderDetail.MenuID
		newOrderDetail.MenuOptionID = incomingOrderDetail.MenuOptionID
		newOrderDetail.MenuOptionValueID = incomingOrderDetail.MenuOptionValueID
		newOrderDetail.Quantity = incomingOrderDetail.Quantity

		newOrderDetails = append(newOrderDetails, newOrderDetail)
	}

	result = database.DB.Create(newOrderDetails)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when inserting order details",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "new order successfully created",
	})
}

func CompleteOrder(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid order ID",
		})
	}

	order := new(models.Order)

	order.ID = orderID
	order.CompletedAt = time.Now()
	result := database.DB.Updates(&order)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "failed to complete order",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "order completed",
	})
}
