package handlers

import (
	"abs-app/database"
	"abs-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	members := new([]models.User)

	result := database.DB.
		Preload("Role").
		Find(&members)

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

func GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "error when parsing member id",
		})
	}

	member := new(models.User)
	result := database.DB.
		Preload("Role").
		Limit(1).
		Find(&member, "id = ?", id)

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
		"data": member,
	})
}

func CreateNewUser(c *fiber.Ctx) error {
	newUser := new(models.User)

	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "Something wrong with your data",
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
		"message": "member successfully created",
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
			"message": "error when parsing member id",
		})
	}

	member := new(models.User)
	result := database.DB.Where("id = ?", id).Delete(member)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "member successfully deleted",
	})
}
