package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWTToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("Shoxrux1801$"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *handlerV1) CreateSuperAdmin(ctx *gin.Context) {
	var req models.CreateSuperAdmin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	existingSuperAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx, req.Username)
	if err == nil && existingSuperAdmin != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "This username already exists",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	superAdmin, err := h.strg.SuperAdmin().Create(ctx, &repo.CreateSuperAdmin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error :(",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateSuperAdmin{
		Username:  superAdmin.Username,
		CreatedAt: superAdmin.CreatedAt,
	})
}

func (h *handlerV1) Login(ctx *gin.Context) {
	var req models.LoginAdmin

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var role string
	var hashedPassword string
	var username string
	var adminId int
	// 1. SuperAdmin jadvalidan tekshirish
	superAdmin, err := h.strg.SuperAdmin().GetByUsername(context.TODO(), req.Username)
	if err == nil && superAdmin != nil {
		role = "super_admin"
		hashedPassword = superAdmin.Password
		username = superAdmin.Username
	} else {
		// 2. Admin jadvalidan tekshirish
		admin, err := h.strg.Admin().GetByUsername(context.TODO(), req.Username)
		if err == nil && admin != nil {
			role = "admin"
			hashedPassword = admin.PasswordHash
			username = admin.Username
			adminId = admin.Id
		} else {
			// 3. Ikkalasida ham topilmasa, xato qaytarish
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
	}

	// 4. Parolni tekshirish (faqat topilgan jadvalga tegishli foydalanuvchining parolini tekshiramiz)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	// 5. JWT token yaratish
	token, err := CreateJWTToken(username, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// 6. Tokenni bazaga saqlash (qaysi jadvaldan topilganiga qarab)
	if role == "super_admin" {
		err = h.strg.SuperAdmin().UpdateToken(ctx, username, token)
	} else {
		admn, err := h.strg.Admin().GetByUsername(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete old token",
			})
			return
		}
		if admn.Status == "inactive" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Your account is blocked",
			})
			return
		}
		tokens, err := h.strg.Token().GetByAdminId(ctx, adminId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch tokens",
			})
			return
		}
		// Agar tokenlar 5 tadan ko‘p bo‘lsa, eng eski tokenni o‘chirish
		if len(tokens) >= 5 {
			// Eng eski tokenni o‘chirish (tokens[0] eng eski bo‘ladi)
			err := h.strg.Token().Delete(ctx, tokens[0].Id)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to delete old token",
				})
				return
			}
		}
		_, err = h.strg.Token().Create(ctx, &repo.CreateToken{
			AdminId: adminId,
			Token:   token,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create token",
			})
			return
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update token in database",
		})
		return
	}

	// 7. Tokenni qaytarish
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  role,
	})
}
