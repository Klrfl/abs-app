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
	menuRoute := app.Group("/menu")

	menuRoute.Get("/", handlers.GetMenu)
	menuRoute.Get("/:id", handlers.GetMenuItemByID)
	menuRoute.Post("/", handlers.CreateNewMenuItem)
	menuRoute.Patch("/:id", handlers.UpdateMenuItem)
	menuRoute.Delete("/:id", handlers.DeleteMenuItem)

	memberRoute := app.Group("/members")
	memberRoute.Get("/", handlers.GetMembers)
	memberRoute.Get("/:id", handlers.GetMemberByID)
	memberRoute.Post("/", handlers.CreateNewMember)
	memberRoute.Patch("/:id", handlers.UpdateMemberData)
	memberRoute.Delete("/:id", handlers.DeleteMember)

	ordersRoute := app.Group("/orders")
	ordersRoute.Get("/", handlers.GetOrders)
	ordersRoute.Get("/:id", handlers.GetOrderByID)
	ordersRoute.Post("/", handlers.CreateNewOrder)
	ordersRoute.Patch("/:id", handlers.CompleteOrder)
}
