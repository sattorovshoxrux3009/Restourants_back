package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = []byte("Shoxrux1801$")

func (h *handlerV1) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization error"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization error"})
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token error"})
		}

		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token error"})
		}

		c.Locals("username", username) // username saqlash
		return c.Next()
	}
}
func (h *handlerV1) SuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("username")
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		superAdminToken, err := h.strg.SuperAdmin().GetToken(c.Context(), username.(string))
		if err != nil || superAdminToken != tokenString {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "SuperAdmin authorization failed"})
		}

		c.Locals("sadmin_username", username)
		c.Locals("role", "super_admin")
		return c.Next()
	}
}
func (h *handlerV1) AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		username := c.Locals("username")

		admin, err := h.strg.Admin().GetByUsername(c.Context(), username.(string))
		if err != nil || admin == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin not found"})
		}

		adminTokens, err := h.strg.Token().GetByAdminId(c.Context(), admin.Id)
		if err != nil || len(adminTokens) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Admin token not found"})
		}

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

// package v1

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// )

// // JWT ni yaratishda ishlatilgan secret key
// var SecretKey = []byte("Shoxrux1801$")

// func (h *handlerV1) AuthMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if authHeader == "" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
// 			ctx.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == authHeader {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
// 			ctx.Abort()
// 			return
// 		}

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
// 			return SecretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
// 			ctx.Abort()
// 			return
// 		}

// 		// Token ichidagi username ni saqlaymiz
// 		username, ok := claims["username"].(string)
// 		if !ok {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
// 			ctx.Abort()
// 			return
// 		}
// 		ctx.Set("username", username) // username saqlash
// 		ctx.Next()
// 	}
// }

// func (h *handlerV1) SuperAdminMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		username, _ := c.Get("username")
// 		authHeader := c.GetHeader("Authorization")
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		superAdminToken, err := h.strg.SuperAdmin().GetToken(c, username.(string))
// 		if err != nil || superAdminToken != tokenString {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "SuperAdmin authorization failed"})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("sadmin_username", username)
// 		c.Set("role", "super_admin")
// 		c.Next()
// 	}
// }

// func (h *handlerV1) AdminMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		username, _ := c.Get("username")
// 		admin, err := h.strg.Admin().GetByUsername(c, username.(string))
// 		if err != nil || admin == nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
// 			c.Abort()
// 			return
// 		}
// 		adminTokens, err := h.strg.Token().GetByAdminId(c, admin.Id)
// 		if err != nil || len(adminTokens) == 0 {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin token not found"})
// 			c.Abort()

// 			return
// 		}
// 		validToken := false
// 		for _, token := range adminTokens {
// 			if token.Token == tokenString {
// 				validToken = true
// 				break
// 			}
// 		}
// 		if !validToken {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("admin_id", admin.Id)
// 		c.Set("role", "admin")
// 		c.Next()
// 	}
// }
