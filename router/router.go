package router

import (
	"abs-app/handlers"
	"abs-app/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
	}))

	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)
	app.Post("/signup", handlers.Signup)

	api := app.Group("/api")

	menuRoute := api.Group("/menu")
	menuRoute.Get("/", handlers.GetMenu)
	menuRoute.Get("/:id", handlers.GetMenuItemByID)

	api.Post("/orders", handlers.CreateNewAnonymousOrder)

	api.Use(middleware.ValidateUserJWT)

	userRoute := api.Group("/user")
	userRoute.Get("/", handlers.GetUserForUser)
	userRoute.Patch("/", handlers.UpdateUserDataForUser)

	ordersRoute := api.Group("/orders")
	ordersRoute.Get("/", handlers.GetOrdersForUser)
	ordersRoute.Get("/:id", handlers.GetOrdersForUserByID)
	ordersRoute.Post("/", handlers.CreateNewOrder)

	// admin routes
	adminRoute := app.Group("/api/admin", middleware.ValidateAdminJWT)

	adminUsersRoute := adminRoute.Group("/users")
	adminUsersRoute.Get("/", handlers.GetUsers)
	adminUsersRoute.Post("/", handlers.CreateNewUser)
	adminUsersRoute.Get("/:id", handlers.GetUserByID)
	adminUsersRoute.Patch("/:id", handlers.UpdateUserData)
	adminUsersRoute.Delete("/:id", handlers.DeleteUser)

	adminOrdersRoute := adminRoute.Group("/orders")
	adminOrdersRoute.Get("/", handlers.GetOrders)
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
