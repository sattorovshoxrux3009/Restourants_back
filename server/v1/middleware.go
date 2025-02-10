package v1

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT ni yaratishda ishlatilgan secret key
var SecretKey = []byte("Shoxrux1801$")

func (h *handlerV1) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
			ctx.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
			ctx.Abort()
			return
		}

		// Token ichidagi username ni saqlaymiz
		username, ok := claims["username"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
			ctx.Abort()
			return
		}
		ctx.Set("username", username) // username saqlash
		ctx.Next()
	}
}

func (h *handlerV1) SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		superAdminToken, err := h.strg.SuperAdmin().GetToken(c, username.(string))
		if err != nil || superAdminToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "SuperAdmin authorization failed"})
			c.Abort()
			return
		}
		c.Set("role", "super_admin")
		c.Next()
	}
}

func (h *handlerV1) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		username, _ := c.Get("username")
		admin, err := h.strg.Admin().GetByUsername(c, username.(string))
		if err != nil || admin == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
			c.Abort()
			return
		}
		adminTokens, err := h.strg.Token().GetByAdminId(c, admin.Id)
		if err != nil || len(adminTokens) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin token not found"})
			c.Abort()

			return
		}
		validToken := false
		for _, token := range adminTokens {
			if token.Token == tokenString {
				validToken = true
				break
			}
		}
		if !validToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("role", "admin")
		c.Next()
	}
}

// func (h *handlerV1) AuthMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// Headerdan tokenni olish
// 		authHeader := ctx.GetHeader("Authorization")
// 		if authHeader == "" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
// 			ctx.Abort()
// 			return
// 		}

// 		// "Bearer " qismi bor yoki yo‘qligini tekshiramiz
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == authHeader { // Agar token "Bearer " bilan boshlanmasa
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
// 			ctx.Abort()
// 			return
// 		}

// 		// Tokenni tekshirish
// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
// 			return SecretKey, nil
// 		})

// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
// 			ctx.Abort()
// 			return
// 		}

// 		// Token yaroqsizligi tekshirildi
// 		if !token.Valid {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
// 			ctx.Abort()
// 			return
// 		}

// 		username, ok := claims["username"].(string)
// 		if !ok {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
// 			ctx.Abort()
// 			return
// 		}

// 		// SuperAdmin bazasidan tokenni olish
// 		superAdminToken, err := h.strg.SuperAdmin().GetToken(ctx, username)
// 		if err != nil || superAdminToken != tokenString {
// 			// Admin bazasidan tokenni olish
// 			// adminToken, err := h.strg.Admin().GetToken(ctx, username)
// 			// if err != nil || adminToken != tokenString {
// 			// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
// 			// 	ctx.Abort()
// 			// 	return
// 			// }
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
// 			ctx.Abort()
// 			return
// 			// Admin bo‘lsa, rolni 'admin' deb belgilaymiz
// 			//ctx.Set("role", "admin")
// 		} else {
// 			// SuperAdmin bo‘lsa, rolni 'super_admin' deb belgilaymiz
// 			ctx.Set("role", "super_admin")
// 		}

// 		// Token ichidagi `username`, `user_id` ni saqlaymiz
// 		ctx.Set("username", claims["username"]) // username

// 		// So‘rovni davom ettiramiz
// 		ctx.Next()
// 	}
// }
