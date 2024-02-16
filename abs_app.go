package main

import (
	"abs-app/database"
	"abs-app/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Init()

	app := fiber.New()
	router.SetupRoutes(app)

	app.Listen(":8080")
}
