package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,              // Maksimal 60 ta so'rov
		Expiration: 1 * time.Minute, // 1 daqiqa ichida
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // IP manzil bo'yicha cheklash
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "So'rovlar soni cheklangan, keyinroq urinib ko'ring",
			})
		},
	})
}
