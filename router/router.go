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
	menuRoute.Post("/", handlers.CreateNewDrink)
	menuRoute.Patch("/:id", handlers.UpdateDrink)
	menuRoute.Delete("/:id", handlers.DeleteDrink)

	memberRoute := app.Group("/members")
	memberRoute.Get("/", handlers.GetAllMembers)
	memberRoute.Get("/:id", handlers.GetMemberByID)
	memberRoute.Post("/", handlers.CreateNewMember)
	memberRoute.Patch("/:id", handlers.UpdateMemberData)
	memberRoute.Delete("/:id", handlers.DeleteMember)

	ordersRoute := app.Group("/orders")
	ordersRoute.Get("/", handlers.GetAllOrders)
}
