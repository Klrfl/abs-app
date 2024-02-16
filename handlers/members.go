package handlers

import (
	"abs-app/database"
	"abs-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	users := new([]models.User)

	result := database.DB.
		Preload("Role").
		Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database for users",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": users,
		})
	}

	return c.JSON(fiber.Map{
		"err":  false,
		"data": users,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "ID not valid",
		})
	}

	user := new(models.User)
	result := database.DB.
		Preload("Role").
		Limit(1).
		Find(&user, "id = ?", id)

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
		"data": user,
	})
}

func CreateNewUser(c *fiber.Ctx) error {
	newUser := new(models.User)

	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "Something wrong with user payload",
		})
	}

	result := database.DB.Save(&newUser)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     false,
			"message": "error when querying database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "user successfully created",
	})
}

func UpdateUserData(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"err":     false,
		"message": "not implemented yet",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing user ID",
		})
	}

	user := new(models.User)
	result := database.DB.Where("id = ?", id).Delete(user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "user successfully deleted",
	})
}
