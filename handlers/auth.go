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
	// receive members struct
	var newMember models.Member

	if err := c.BodyParser(&newMember); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the payload",
		})
	}

	// insert to members table and generate generatedPassword with bcrypt
	generatedPassword, err := bcrypt.GenerateFromPassword([]byte(newMember.Password), 14)

	if err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "error when signing up user",
		})
	}

	newMember.Password = string(generatedPassword)
	result := database.DB.Create(&newMember)

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
	var incomingMember models.Member

	if err := c.BodyParser(&incomingMember); err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "something wrong with the payload",
		})
	}

	var existingMember models.Member

	result := database.DB.
		Where("email = ?", incomingMember.Email).
		Find(&existingMember)

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

	if err := bcrypt.CompareHashAndPassword([]byte(existingMember.Password), []byte(incomingMember.Password)); err != nil {
		return c.JSON(fiber.Map{
			"err":     true,
			"message": "wrong email or password",
		})
	}

	expiryTime := time.Now().Add(1 * time.Hour).UTC()
	now := time.Now().UTC()

	claims := &models.JWTClaim{
		Name:  existingMember.Name,
		Email: existingMember.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   existingMember.ID.String(),
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
