package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func IPLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()
		requestTime := time.Now().Format("2006-01-02 15:04:05")
		println("Yangi so'rov! IP:", clientIP, "Vaqt:", requestTime)
		return c.Next()
	}
}
