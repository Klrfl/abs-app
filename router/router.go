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

	app.Use(middleware.ValidateUserJWT)
	menuRoute.Post("/", handlers.CreateNewMenuItem)
	menuRoute.Delete("/", handlers.DeleteMenuItems)
	menuRoute.Patch("/:id", handlers.UpdateMenuItem)
	menuRoute.Post("/:id/variant_values", handlers.InsertNewPrices)
	menuRoute.Patch("/:id/variant_values", handlers.UpdatePrices)
	menuRoute.Delete("/:id/variant_values", handlers.DeletePrice)
	menuRoute.Delete("/:id", handlers.DeleteMenuItem)

	adminRoute := app.Group("/admin")
	adminRoute.Use(middleware.ValidateAdminJWT)

	memberRoute := adminRoute.Group("/users")
	memberRoute.Get("/", handlers.GetUsers)
	memberRoute.Post("/", handlers.CreateNewUser)
	memberRoute.Get("/:id", handlers.GetUserByID)
	memberRoute.Patch("/:id", handlers.UpdateUserData)
	memberRoute.Delete("/:id", handlers.DeleteUser)

	ordersRoute := adminRoute.Group("/orders")
	ordersRoute.Get("/", handlers.GetOrders)
	ordersRoute.Post("/", handlers.CreateNewOrder)
	ordersRoute.Get("/:id", handlers.GetOrderByID)
	ordersRoute.Patch("/:id", handlers.CompleteOrder)
}
