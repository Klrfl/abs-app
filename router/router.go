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

	api := app.Group("/api")

	menuRoute := api.Group("/menu")
	menuRoute.Get("/", handlers.GetMenu)
	menuRoute.Get("/:id", handlers.GetMenuItemByID)

	api.Use(middleware.ValidateUserJWT)

	// TODO: make separate handlers for getting orders for user only
	// or modify handler
	ordersRoute := api.Group("/orders")
	ordersRoute.Get("/", handlers.GetPendingOrders)
	ordersRoute.Get("/:id", handlers.GetOrderByID)
	ordersRoute.Post("/", handlers.CreateNewOrder)

	// admin routes
	adminRoute := app.Group("/admin", middleware.ValidateAdminJWT)

	memberRoute := adminRoute.Group("/users")
	memberRoute.Get("/", handlers.GetUsers)
	memberRoute.Post("/", handlers.CreateNewUser)
	memberRoute.Get("/:id", handlers.GetUserByID)
	memberRoute.Patch("/:id", handlers.UpdateUserData)
	memberRoute.Delete("/:id", handlers.DeleteUser)

	adminOrdersRoute := adminRoute.Group("/orders")
	adminOrdersRoute.Get("/", handlers.GetPendingOrders)
	adminOrdersRoute.Get("/:id", handlers.GetOrderByID)
	adminOrdersRoute.Patch("/:id", handlers.CompleteOrder)

	adminMenuRoute := adminRoute.Group("/menu")

	adminMenuRoute.Post("/", handlers.CreateNewMenuItem)
	adminMenuRoute.Delete("/", handlers.DeleteMenuItems)
	adminMenuRoute.Patch("/:id", handlers.UpdateMenuItem)
	adminMenuRoute.Post("/:id/variant_values", handlers.InsertNewPrices)
	adminMenuRoute.Patch("/:id/variant_values", handlers.UpdatePrices)
	adminMenuRoute.Delete("/:id/variant_values", handlers.DeletePrice)
	adminMenuRoute.Delete("/:id", handlers.DeleteMenuItem)
}
