package main

import (
	"abs-app/database"
	"abs-app/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db := database.Init()

	app := fiber.New()
	router.SetupRoutes(app, db)

	app.Listen(":8080")
}
