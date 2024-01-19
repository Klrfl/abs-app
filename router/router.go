package router

import (
	"abs-app/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New())
	menuRoute := app.Group("/drinks")

	menuRoute.Get("/", handlers.GetAllDrinks)
	menuRoute.Get("/:id", handlers.GetDrinkByID)

	memberRoute := app.Group("/members")
	memberRoute.Get("/", handlers.GetAllMembers)
}
