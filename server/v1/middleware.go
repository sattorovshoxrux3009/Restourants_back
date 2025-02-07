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
		// Headerdan tokenni olish
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
			ctx.Abort()
			return
		}

		// "Bearer " qismi bor yoki yo‘qligini tekshiramiz
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Agar token "Bearer " bilan boshlanmasa
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
			ctx.Abort()
			return
		}

		// Tokenni tekshirish
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
			ctx.Abort()
			return
		}

		// Token yaroqsizligi tekshirildi
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			ctx.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error"})
			ctx.Abort()
			return
		}

		// SuperAdmin bazasidan tokenni olish
		superAdminToken, err := h.strg.SuperAdmin().GetToken(ctx, username)
		if err != nil || superAdminToken != tokenString {
			// Admin bazasidan tokenni olish
			// adminToken, err := h.strg.Admin().GetToken(ctx, username)
			// if err != nil || adminToken != tokenString {
			// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			// 	ctx.Abort()
			// 	return
			// }
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			ctx.Abort()
			return
			// Admin bo‘lsa, rolni 'admin' deb belgilaymiz
			//ctx.Set("role", "admin")
		} else {
			// SuperAdmin bo‘lsa, rolni 'super_admin' deb belgilaymiz
			ctx.Set("role", "super_admin")
		}

		// Token ichidagi `username`, `user_id` ni saqlaymiz
		ctx.Set("username", claims["username"]) // username

		// So‘rovni davom ettiramiz
		ctx.Next()
	}
}
