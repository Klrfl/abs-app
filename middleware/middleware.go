package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func decodeToken(tokenString string, key string) (*jwt.Token, error) {
	decodedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %s", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	return decodedToken, nil
}

func getTokenString(c *fiber.Ctx) string {
	if len(c.Get("Authorization")) != 0 {
		return strings.Split(c.Get("Authorization"), " ")[1]
	}
	return c.Cookies("token")
}

func ValidateUserJWT(c *fiber.Ctx) error {
	//TODO: verify tokenString and refresh if less then 20 minutes
	tokenString := getTokenString(c)

	key := os.Getenv("SECRET")

	decodedToken, err := decodeToken(tokenString, key)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err":     true,
			"message": "cannot verify token",
		})
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		c.Locals("user_id", claims["ID"])
		expiryTime := claims["exp"].(float64)

		// TODO: refresh token
		if int64(expiryTime) < time.Now().Unix() {
			c.Cookie(&fiber.Cookie{
				Name:  "token",
				Value: "",
			})

			return c.SendStatus(fiber.StatusUnauthorized)
		}
	} else {
		log.Println("error when parsing claims")
	}

	return c.Next()
}

func ValidateAdminJWT(c *fiber.Ctx) error {
	tokenString := getTokenString(c)
	key := os.Getenv("SECRET")

	decodedToken, err := decodeToken(tokenString, key)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"err":     true,
			"message": "unable to verify token",
		})
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		if claims["RoleID"].(float64) == 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"err":     true,
				"message": "only users with permission can access this resource",
			})
		}
	}
	return c.Next()
}
