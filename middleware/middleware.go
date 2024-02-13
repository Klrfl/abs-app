package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(c *fiber.Ctx) error {
	// verify tokenString and refresh if less then 20 minutes
	tokenString := c.Cookies("token")
	key := os.Getenv("SECRET")

	decodedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %s", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err":     true,
			"message": "cannot verify token",
		})
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		expiryTime := claims["exp"].(float64)

		// refresh token or abort request
		if int64(expiryTime) < time.Now().Unix() {
			c.Cookie(&fiber.Cookie{
				Name:  "token",
				Value: "",
			})

			return c.SendStatus(fiber.StatusUnauthorized)
		}
		log.Println(claims)
	} else {
		log.Println("error when parsing claims")
	}

	return c.Next()
}
