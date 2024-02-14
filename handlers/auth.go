package handlers

import (
	"abs-app/database"
	"abs-app/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	var newUser models.User

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the payload",
		})
	}

	generatedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when signing up user",
		})
	}

	newUser.Password = string(generatedPassword)
	newUser.RoleID = 1
	result := database.DB.Create(&newUser)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when signing up user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"err":     false,
		"message": "signup success! redirect user to login page",
	})
}

func Login(c *fiber.Ctx) error {
	//TODO: verify password and email
	var incomingUser models.User

	if err := c.BodyParser(&incomingUser); err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the payload",
		})
	}

	var existingUser models.User

	result := database.DB.
		Where("email = ?", incomingUser.Email).
		Limit(1).
		Find(&existingUser)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when verifying user",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err":     true,
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(incomingUser.Password)); err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "wrong email or password",
		})
	}

	expiryTime := time.Now().Add(1 * time.Hour).UTC()
	now := time.Now().UTC()

	claims := &models.JWTClaim{
		Name:   existingUser.Name,
		Email:  existingUser.Email,
		RoleID: existingUser.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   existingUser.ID.String(),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := os.Getenv("SECRET")
	signedToken, err := newToken.SignedString([]byte(key))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err":     true,
			"message": "error when signing token",
		})
	}

	// then set a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    signedToken,
		Path:     "/",
		Expires:  expiryTime,
		MaxAge:   int(expiryTime.Unix()),
		Secure:   true,
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "successfully logged in",
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})

	return c.JSON(fiber.Map{
		"err":     false,
		"message": "successfully logged out",
	})
}
