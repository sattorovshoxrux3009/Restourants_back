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

	router.POST("/v1/create-s-admin", handler.CreateSuperAdmin)
	router.POST("/v1/login", handler.Login)

	// router.POST("/v1/login", handler.LoginUser)
	// // Himoyalangan API-lar (Middleware ishlaydi)
	protected := router.Group("/v1")
	protected.Use(handler.AuthMiddleware()) // JWT tokenni tekshirish
	{
		protected.POST("/create-admin", handler.CreateAdmin) 
		protected.GET("/all-admins",handler.GetAllAdmins)
	}
	return router
}
