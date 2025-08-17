package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sattorovshoxrux3009/Restourants_back/middleware"
	v1 "github.com/sattorovshoxrux3009/Restourants_back/server/v1"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
)

type Options struct {
	Strg storage.StorageI
}

func NewServer(opts *Options) *fiber.App {
	app := fiber.New()

	// Global middlewarelar
	app.Use(middleware.IPLogger())       // IP logger middleware
	app.Use(middleware.CorsMiddleware()) // CORS middleware
	app.Use(middleware.RateLimiter())    // Rate limiter middleware

	// Swagger documentation endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Handler
	handler := v1.New(&v1.HandlerV1{
		Strg: opts.Strg,
	})

	// Middleware'larni initializatsiya qilish
	authMiddleware := middleware.AuthMiddleware()
	adminMiddleware := middleware.NewAdminMiddleware(opts.Strg)
	superAdminMiddleware := middleware.NewSuperAdminMiddleware(opts.Strg)

	// OPTIONS so'rovlari uchun
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// Statik fayllar uchun
	app.Static("/uploads", "./uploads")

	// Public endpoints
	app.Post("/v1/login", handler.Login)

	// Super Admin routes
	superAdmin := app.Group("/v1/superadmin",
		authMiddleware,
		superAdminMiddleware.SuperAdminAuth(),
	)
	{
		superAdmin.Get("/admins/:id?", handler.GetAdmins)
		superAdmin.Get("/admins/:id/details", handler.GetAdminDetails)
		superAdmin.Get("/restaurants/:id?", handler.GetSRestaurants)
		superAdmin.Get("/restaurants/:id/details", handler.GetSRestaurantDetails)
		superAdmin.Get("/menu/:id?", handler.GetSMenu)
		superAdmin.Get("/profile", handler.GetSProfile)
		superAdmin.Get("/event-prices/:id?", handler.GetSEventPrices)

		superAdmin.Post("/admin", handler.CreateAdmin)
		superAdmin.Post("/restaurant", handler.CreateSRestaurant)
		superAdmin.Post("/menu", handler.CreateSMenu)
		superAdmin.Post("/event-prices", handler.CreateSEventPrices)

		superAdmin.Put("/admin/:id", handler.UpdateAdmin)
		superAdmin.Put("/restaurants/:id/status", handler.UpdateSRestaurantStatus)
		superAdmin.Put("/restaurants/:id", handler.UpdateSRestaurant)
		superAdmin.Put("/menu/:id", handler.UpdateSMenu)
		superAdmin.Put("/profile", handler.UpdateSProfile)
		superAdmin.Put("/event-prices/:id", handler.UpdateSEventPrices)

		superAdmin.Delete("/admin/:id", handler.DeleteAdmin)
		superAdmin.Delete("/restaurants/:id", handler.DeleteSRestaurant)
		superAdmin.Delete("/menu/:id", handler.DeleteSMenu)
		superAdmin.Delete("/event-prices/:id", handler.DeleteSEventPrices)
	}
	// Admin routes
	admin := app.Group("/v1/admin",
		authMiddleware,
		adminMiddleware.AdminAuth(),
	)
	{
		admin.Get("/profile", handler.GetProfile)
		admin.Get("/restaurants/:id?", handler.GetARestaurants)
		admin.Get("/menu/:id?", handler.GetAMenu)
		admin.Get("/event-prices/:id?", handler.GetAEventPrices)

		admin.Put("/profile", handler.UpdateProfile)
		admin.Put("/restaurants/:id", handler.UpdateARestauranats)
		admin.Put("/menu/:id", handler.UpdateAMenu)
		admin.Put("/event-prices/:id", handler.UpdateAEventPrices)

		admin.Post("/restaurants", handler.CreateARestaurant)
		admin.Post("/menu", handler.CreateAMenu)
		admin.Post("/event-prices", handler.CreateAEventPrices)

		admin.Delete("/restaurants/:id", handler.DeleteARestaurants)
		admin.Delete("/menu/:id", handler.DeleteAMenu)
		admin.Delete("/event-prices/:id", handler.DeleteAEventPrices)
	}

	app.Get("/v1/restaurants/:id?", handler.GetRestaurants)
	app.Get("/v1/menu/:id?", handler.GetMenu)
	app.Get("/v1/restaurants/:id/details", handler.GetRestaurantDetails)
	app.Get("/v1/search", handler.GetSearch)

	return app
}
