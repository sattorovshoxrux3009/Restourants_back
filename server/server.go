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
	app := fiber.New()

	// IP log middleware
	// app.Use(func(c *fiber.Ctx) error {
	// 	clientIP := c.IP()
	// 	println("Yangi so‘rov! IP:", clientIP)
	// 	return c.Next()
	// })

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
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		AllowCredentials: true,
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

	// Statik fayllar uchun
	app.Static("/uploads", "./uploads")

	// Auth yo‘nalishlari
	app.Post("/v1/create-s-admin", handler.CreateSuperAdmin)
	app.Post("/v1/login", handler.Login)

	// Restoran va menyu yo‘nalishlari
	app.Get("/v1/restaurants", handler.GetRestourants)
	app.Get("/v1/restaurants/:id", handler.GetRestourants)
	app.Get("/v1/menu", handler.GetMenu)
	app.Get("/v1/menu/:id", handler.GetMenu)

	// Super Admin yo‘nalishlari
	superAdmin := app.Group("/v1", handler.AuthMiddleware(), handler.SuperAdminMiddleware())
	{
		superAdmin.Get("/admins", handler.GetAdmins)
		superAdmin.Get("/admins/:id", handler.GetAdmins)
		superAdmin.Get("/admins/:id/details", handler.GetAdminDetails)
		superAdmin.Get("/s-restaurants", handler.GetSRestourants)
		superAdmin.Get("/s-restaurants/:id", handler.GetSRestourants)
		superAdmin.Get("/s-restaurants/:id/details", handler.GetSRestaurantDetails)
		superAdmin.Get("/s-menu", handler.GetSMenu)
		superAdmin.Get("/s-menu/:id", handler.GetSMenu)
		superAdmin.Get("/s-profile", handler.GetSProfile)
		superAdmin.Get("/s-event-prices", handler.GetSEventPrices)
		superAdmin.Get("/s-event-prices/:id", handler.GetSEventPrices)

		superAdmin.Post("/create-admin", handler.CreateAdmin)
		superAdmin.Post("/create-restaurant", handler.CreateRestaurant)
		superAdmin.Post("/s-menu", handler.CreateSMenu)
		superAdmin.Post("/s-event-prices", handler.CreateSEventPrices)

		superAdmin.Put("/update-admin/:id", handler.UpdateAdmin)
		superAdmin.Put("/restaurants/:id/status", handler.UpdateRestaurantStatus)
		superAdmin.Put("/restaurants/:id", handler.UpdateRestaurant)
		superAdmin.Put("/s-menu/:id", handler.UpdateSMenu)
		superAdmin.Put("/s-profile", handler.UpdateSProfile)
		superAdmin.Put("/s-event-prices/:id", handler.UpdateSEventPrices)

		superAdmin.Delete("/admin/:id", handler.DeleteAdmin)
		superAdmin.Delete("/s-restaurants/:id", handler.DeleteRestaurant)
		superAdmin.Delete("/s-menu/:id", handler.DeleteSMenu)
		superAdmin.Delete("/s-event-prices/:id", handler.DeleteSEventPrices)
	}

	// // Admin yo‘nalishlari
	admin := app.Group("/v1", handler.AuthMiddleware(), handler.AdminMiddleware())
	{
		admin.Get("/profile", handler.GetProfile)
	}

	return app
}

// package server

// import (
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	v1 "github.com/sattorovshoxrux3009/Restourants_back/server/v1"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage"
// )

// type Options struct {
// 	Strg storage.StorageI
// }

// func NewServer(opts *Options) *gin.Engine {
// 	router := gin.New()
// 	router.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:5173"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		AllowCredentials: true,
// 	}))
// 	handler := v1.New(&v1.HandlerV1{
// 		Strg: opts.Strg,
// 	})
// 	router.Static("/uploads", "./uploads")

// 	router.POST("/v1/create-s-admin", handler.CreateSuperAdmin)
// 	router.POST("/v1/login", handler.Login)

// 	router.GET("/v1/restaurants", handler.GetRestourants)
// 	router.GET("/v1/restaurants/:id", handler.GetRestourants)
// 	router.GET("/v1/menu", handler.GetMenu)
// 	router.GET("/v1/menu/:id", handler.GetMenu)

// 	superAdmin := router.Group("/v1")
// 	superAdmin.Use(handler.AuthMiddleware(), handler.SuperAdminMiddleware())
// 	{
// 		superAdmin.GET("/admins", handler.GetAdmins)
// 		superAdmin.GET("/admins/:id", handler.GetAdmins)
// 		superAdmin.GET("/admins/:id/details", handler.GetAdminDetails)
// 		superAdmin.GET("/s-restaurants", handler.GetSRestourants)
// 		superAdmin.GET("/s-restaurants/:id", handler.GetSRestourants)
// 		superAdmin.GET("/s-restaurants/:id/details", handler.GetSRestaurantDetails)
// 		superAdmin.GET("/s-menu", handler.GetSMenu)
// 		superAdmin.GET("/s-menu/:id", handler.GetSMenu)
// 		superAdmin.GET("/s-profile", handler.GetSProfile)
// 		superAdmin.GET("/s-event-prices", handler.GetSEventPrices)
// 		superAdmin.GET("/s-event-prices/:id", handler.GetSEventPrices)

// 		superAdmin.POST("/create-admin", handler.CreateAdmin)
// 		superAdmin.POST("/create-restaurant", handler.CreateRestaurant)
// 		superAdmin.POST("/s-menu", handler.CreateSMenu)
// 		superAdmin.POST("/s-event-prices", handler.CreateSEventPrices)

// 		superAdmin.PUT("/update-admin/:id", handler.UpdateAdmin)
// 		superAdmin.PUT("/restaurants/:id/status", handler.UpdateRestaurantStatus)
// 		superAdmin.PUT("/restaurants/:id", handler.UpdateRestaurant)
// 		superAdmin.PUT("/s-menu/:id", handler.UpdateSMenu)
// 		superAdmin.PUT("/s-profile", handler.UpdateSProfile)
// 		superAdmin.PUT("/s-event-prices/:id", handler.UpdateSEventPrices)

// 		superAdmin.DELETE("/admin/:id", handler.DeleteAdmin)
// 		superAdmin.DELETE("/s-restaurants/:id", handler.DeleteRastaurant)
// 		superAdmin.DELETE("/s-menu/:id", handler.DeleteSMenu)
// 		superAdmin.DELETE("/s-event-prices/:id", handler.DeleteSEventPrices)
// 	}

// 	admin := router.Group("/v1")
// 	admin.Use(handler.AuthMiddleware(), handler.AdminMiddleware())
// 	{
// 		admin.GET("/profile", handler.GetProfile)
// 	}
// 	return router
// }
