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
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
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
	existingSuperAdmin, err := h.strg.SuperAdmin().GetByUserneme(ctx, req.Username)
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
		Username: req.Username,
		Password: string(hashedPassword),
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
	var role string = "superadmin"
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}
	superAdmin, err := h.strg.SuperAdmin().GetByUserneme(context.TODO(), req.Username)
	if err != nil || superAdmin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Parolni hash bilan solishtirish
	err = bcrypt.CompareHashAndPassword([]byte(superAdmin.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password!",
		})
		return
	}

	// JWT token yaratish
	token, err := CreateJWTToken(req.Username, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Tokenni bazaga saqlash
	err = h.strg.SuperAdmin().UpdateToken(ctx, req.Username, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update token in database",
		})
		return
	}

	// Tokenni qaytarish
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  role,
	})
}
