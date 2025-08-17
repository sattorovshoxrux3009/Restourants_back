package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
)

type SuperAdminMiddleware struct {
	strg storage.StorageI
}

func NewSuperAdminMiddleware(strg storage.StorageI) *SuperAdminMiddleware {
	return &SuperAdminMiddleware{
		strg: strg,
	}
}

func (s *SuperAdminMiddleware) SuperAdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("username")
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		superAdminToken, err := s.strg.SuperAdmin().GetToken(c.Context(), username.(string))
		if err != nil || superAdminToken != tokenString {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "SuperAdmin authorization failed"})
		}

		c.Locals("sadmin_username", username)
		c.Locals("role", "super_admin")
		return c.Next()
	}
}
