package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	v1 "github.com/sattorovshoxrux3009/Restourants_back/server/v1"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
)

type Options struct {
	Strg storage.StorageI
}

func NewServer(opts *Options) *fiber.App {
	// app := fiber.New(fiber.Config{
	// 	// HTTPS uchun TLS konfiguratsiyasi
	// 	DisableKeepalive: false,
	// })
	app := fiber.New()
	// IP log middleware
	app.Use(func(c *fiber.Ctx) error {
		clientIP := c.IP()
		requestTime := time.Now().Format("2006-01-02 15:04:05") // Yil-oy-kun soat:minut:sekund
		println("Yangi so‘rov! IP:", clientIP, "Vaqt:", requestTime)
		return c.Next()
	})

	// var blockedIPs = map[string]bool{
	// 	"172.25.25.101": true, // Bloklangan IP
	// 	"10.10.10.5":    true, // Bloklangan IP
	// }

	// app.Use(func(c *fiber.Ctx) error {
	// 	clientIP := c.IP()

	// 	if blockedIPs[clientIP] {
	// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 			"error": "Sizning IP-manzilingiz bloklangan",
	// 		})
	// 	}

	// 	println("Yangi so‘rov! IP:", clientIP)
	// 	return c.Next()
	// })

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		AllowOriginsFunc: func(origin string) bool { return true }, // OPTIONS muammosini hal qiladi
		// AllowCredentials: true,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        60,              // Maksimal 100 ta so‘rov
		Expiration: 1 * time.Minute, // 1 daqiqa ichida hisoblash
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Foydalanuvchi IP manzili bo‘yicha cheklash
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "So‘rovlar soni cheklangan, keyinroq urinib ko‘ring.",
			})
		},
	}))

	// Handler
	handler := v1.New(&v1.HandlerV1{
		Strg: opts.Strg,
	})
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No Content
	})
	// Statik fayllar uchun
	app.Static("/uploads", "./uploads")

	// app.Post("/v1/s-admin", handler.CreateSuperAdmin)
	app.Post("/v1/login", handler.Login)

	superAdmin := app.Group("/v1/superadmin", handler.AuthMiddleware(), handler.SuperAdminMiddleware())
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
	admin := app.Group("/v1/admin", handler.AuthMiddleware(), handler.AdminMiddleware())
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
