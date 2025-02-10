package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

func (h *handlerV1) CreateAdmin(ctx *gin.Context) {
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

func (h *handlerV1) GetAdmins(ctx *gin.Context) {
	adminID := ctx.Param("id")
	if adminID != "" {
		// ID bo‘yicha bitta adminni olish
		id, err := strconv.Atoi(adminID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
			return
		}

		admin, err := h.strg.Admin().GetById(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
			return
		}

		ctx.JSON(http.StatusOK, admin)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Sahifani 1 dan boshlash
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	
	// Query parametrlarini olish
	status := ctx.Query("status")
	firstname := ctx.Query("firstname")     // firstname qidiruvi
	lastname := ctx.Query("lastname")       // lastname qidiruvi
	email := ctx.Query("email")             // email qidiruvi
	phonenumber := ctx.Query("phonenumber") // phonenumber qidiruvi
	username := ctx.Query("username")       // username qidiruvi

	// Sahifani 1 dan boshlash
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Adminlarni olish
	admins, currentPage, totalPage, err := h.strg.Admin().GetAll(
		ctx,
		status,
		firstname,   // firstname bo'yicha filter
		lastname,    // lastname bo'yicha filter
		email,       // email bo'yicha filter
		phonenumber, // phonenumber bo'yicha filter
		username,    // username bo'yicha filter
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admins"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page":       currentPage,
		"total_page": totalPage,
		"admins":     admins,
	})
}

func (h *handlerV1) UpdateAdminStatus(ctx *gin.Context) {
	adminID := ctx.Param("id") // URL parametri orqali admin ID
	var requestBody models.UpdateAdminStatus

	// JSON ma'lumotlarini bind qilish
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// IDni integerga o'zgartirish
	id, err := strconv.Atoi(adminID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Adminni topish
	admin, err := h.strg.Admin().GetById(ctx, id) // GetById funksiyasini ishlatamiz
	if admin == nil || err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Statusni yangilash
	err = h.strg.Admin().UpdateStatus(ctx, id, requestBody.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
		return
	}

	// Yangilangan statusni va admin ma'lumotlarini yuborish
	ctx.JSON(http.StatusOK, gin.H{"message": "Admin status updated successfully"})
}

func (h *handlerV1) UpdateAdmin(ctx *gin.Context) {
	adminID := ctx.Param("id")
	var requestBody models.UpdateAdmin
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	id, err := strconv.Atoi(adminID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}
	admin, err := h.strg.Admin().GetById(ctx, id) // GetById funksiyasini ishlatamiz
	if admin == nil || err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	var hashedPassword []byte
	if requestBody.Password != "" {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	hashedPasswordStr := string(hashedPassword)

	updateData := &repo.UpdateAdmin{
		FirstName:    requestBody.FirstName,
		LastName:     requestBody.LastName,
		Email:        requestBody.Email,
		PhoneNumber:  requestBody.PhoneNumber,
		Username:     requestBody.Username,
		PasswordHash: hashedPasswordStr,
	}

	// Adminni yangilash
	err = h.strg.Admin().Update(ctx, id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "err"})
		return
	}

	// Yangilangan statusni va admin ma'lumotlarini yuborish
	ctx.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
}

func (h *handlerV1) GetAdminDetails(ctx *gin.Context){
	role, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
		return
	}

	if role != "super_admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update admin status",
		})
		return
	}

}