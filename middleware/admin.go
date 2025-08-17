package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
)

type AdminMiddleware struct {
	strg storage.StorageI
}

func NewAdminMiddleware(strg storage.StorageI) *AdminMiddleware {
	return &AdminMiddleware{
		strg: strg,
	}
}

func (a *AdminMiddleware) AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("username")
		admin, err := a.strg.Admin().GetByUsername(c.Context(), username.(string))
		if err != nil || admin == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
		}

		adminTokens, err := a.strg.Token().GetByAdminId(c.Context(), int(admin.Id))
		if err != nil || len(adminTokens) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin token not found"})
		}

		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		validToken := false
		for _, token := range adminTokens {
			if token.Token == tokenString {
				validToken = true
				break
			}
		}

		if !validToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("admin_id", admin.Id)
		c.Locals("role", "admin")
		return c.Next()
	}
}
