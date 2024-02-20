package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	users := new([]models.User)

	result := database.DB.
		Preload("Role").
		Omit("password").
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
		Omit("password").
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

	tx := database.DB.Begin()
	result := tx.Save(&newUser)

	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     false,
			"message": "error when querying database",
		})
	}

	tx.Commit()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "user successfully created",
	})
}

func UpdateUserData(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "invalid user ID",
		})
	}

	incomingUser, existingUser := new(models.User), new(models.User)

	if err := c.BodyParser(&incomingUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with user payload",
		})
	}

	result := database.DB.
		Where("id = ?", userID).
		Find(&existingUser)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "no user found",
		})
	}

	tx := database.DB.Begin()
	result = tx.
		Where("id = ?", userID).
		Updates(&incomingUser)

	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when updating user",
		})
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     false,
			"message": "error doesn't work",
		})
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"err":     false,
		"message": fmt.Sprintf("user with ID %s successfully updated", userID),
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

	tx := database.DB.Begin()
	result := tx.
		Where("id = ?", id).
		Delete(user)

	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database",
		})
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"err":     false,
		"message": "user successfully deleted",
	})
}

func GetUserForUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"messgae": "user ID not valid",
		})
	}

	existingUser := new(models.User)

	err = database.DB.
		Preload("Role").
		Where("id = ?", userID).
		Omit("password").
		Find(&existingUser).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database for user",
		})
	}

	return c.JSON(fiber.Map{
		"err":     true,
		"message": existingUser,
	})
}

func UpdateUserDataForUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"messgae": "user ID not valid",
		})
	}

	existingUser, incomingUser := new(models.User), new(models.User)

	if err := c.BodyParser(&incomingUser); err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with user payload",
		})
	}

	err = database.DB.
		Preload("Role").
		Where("id = ?", userID).
		Find(&existingUser).Error

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when querying database for user",
		})
	}

	tx := database.DB.Begin()
	err = tx.
		Where("id = ?", userID).
		Updates(&incomingUser).Error

	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when updating user data",
		})
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"err":     true,
		"message": "user data succesfully updated",
	})
}
