package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/sattorovshoxrux3009/Restourants_back/server/v1"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
)

type Options struct {
	Strg storage.StorageI
}

func NewServer(opts *Options) *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	handler := v1.New(&v1.HandlerV1{
		Strg: opts.Strg,
	})
	router.Static("/uploads", "./uploads")

	router.POST("/v1/create-s-admin", handler.CreateSuperAdmin)
	router.POST("/v1/login", handler.Login)

	router.GET("/v1/restaurants", handler.GetRestourants)
	router.GET("/v1/restaurants/:id", handler.GetRestourants)

	superAdmin := router.Group("/v1")
	superAdmin.Use(handler.AuthMiddleware(), handler.SuperAdminMiddleware())
	{
		superAdmin.GET("/admins", handler.GetAdmins)
		superAdmin.GET("/admins/:id", handler.GetAdmins)
		superAdmin.GET("/admins/:id/details", handler.GetAdminDetails)
		superAdmin.GET("/s-restaurants", handler.GetSRestourants)
		superAdmin.GET("/s-restaurants/:id", handler.GetSRestourants)
		superAdmin.POST("/create-admin", handler.CreateAdmin)
		superAdmin.PUT("/update-admin/:id/status", handler.UpdateAdminStatus)
		superAdmin.PUT("/update-admin/:id/limit", handler.UpdateAdminLimit)
		superAdmin.PUT("/update-admin/:id", handler.UpdateAdmin)
		superAdmin.POST("/create-restaurant", handler.CreateRestaurant)
		superAdmin.DELETE("/admin/:id", handler.DeleteAdmin)
	}

	admin := router.Group("/v1")
	admin.Use(handler.AuthMiddleware(), handler.AdminMiddleware())
	{

	}
	return router
}
