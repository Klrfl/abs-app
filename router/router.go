package router

import (
	"abs-app/handlers"
	"abs-app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(limiter.New())

	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)
	app.Post("/signup", handlers.Signup)

	menuRoute := app.Group("/menu")

	menuRoute.Get("/", handlers.GetMenu)
	menuRoute.Get("/:id", handlers.GetMenuItemByID)

	app.Use(middleware.CheckAuth)
	menuRoute.Post("/", handlers.CreateNewMenuItem)
	menuRoute.Delete("/", handlers.DeleteMenuItems)
	menuRoute.Patch("/:id", handlers.UpdateMenuItem)
	menuRoute.Post("/:id/variant_values", handlers.InsertNewPrices)
	menuRoute.Patch("/:id/variant_values", handlers.UpdatePrices)
	menuRoute.Delete("/:id/variant_values", handlers.DeletePrice)
	menuRoute.Delete("/:id", handlers.DeleteMenuItem)

	memberRoute := app.Group("/members")
	memberRoute.Get("/", handlers.GetMembers)
	memberRoute.Post("/", handlers.CreateNewMember)
	memberRoute.Get("/:id", handlers.GetMemberByID)
	memberRoute.Patch("/:id", handlers.UpdateMemberData)
	memberRoute.Delete("/:id", handlers.DeleteMember)

	ordersRoute := app.Group("/orders")
	ordersRoute.Get("/", handlers.GetOrders)
	ordersRoute.Post("/", handlers.CreateNewOrder)
	ordersRoute.Get("/:id", handlers.GetOrderByID)
	ordersRoute.Patch("/:id", handlers.CompleteOrder)
}
