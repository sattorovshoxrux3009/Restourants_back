package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

func (h *handlerV1) CreateAdmin(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error"})
		return
	}
	if role != "super_admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "You do not have permission",
		})
		return
	}

	var req models.CreateAdmin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	existingAdmin, err := h.strg.Admin().GetByUsername(ctx, req.Username)
	if err == nil && existingAdmin != nil {
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

	admin, err := h.strg.Admin().Create(ctx, &repo.CreateAdmin{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error :(",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateAdmin{
		Username: admin.Username,
	})
}

func (h *handlerV1) GetAllAdmins(ctx *gin.Context) {
	// Ro‘lda super_admin borligini tekshirish
	role, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error"})
		return
	}
	if role != "super_admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "You do not have permission",
		})
		return
	}

	// Adminlarni olish
	admins, err := h.strg.Admin().GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admins"})
		return
	}

	// Adminlarni muvaffaqiyatli qaytarish
	ctx.JSON(http.StatusOK, admins)
}
