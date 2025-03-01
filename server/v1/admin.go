package v1

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

func (h *handlerV1) CreateAdmin(ctx *fiber.Ctx) error {
	var req models.CreateAdmin
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingAdmin, err := h.strg.Admin().GetByUsername(ctx.Context(), req.Username)
	if err == nil && existingAdmin != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This username already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	admin, err := h.strg.Admin().Create(ctx.Context(), &repo.CreateAdmin{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error :("})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"username": admin.Username})
}

func (h *handlerV1) GetAdmins(ctx *fiber.Ctx) error {
	adminID := ctx.Params("id")
	if adminID != "" {
		id, err := strconv.Atoi(adminID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid admin ID"})
		}

		admin, err := h.strg.Admin().GetById(ctx.Context(), id)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin not found"})
		}

		return ctx.Status(fiber.StatusOK).JSON(admin)
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	admins, currentPage, totalPage, err := h.strg.Admin().GetAll(
		ctx.Context(),
		ctx.Query("status"),
		ctx.Query("firstname"),
		ctx.Query("lastname"),
		ctx.Query("email"),
		ctx.Query("phonenumber"),
		ctx.Query("username"),
		page,
		limit,
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get admins"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":       currentPage,
		"total_page": totalPage,
		"admins":     admins,
	})
}

func (h *handlerV1) UpdateAdmin(ctx *fiber.Ctx) error {
	adminID := ctx.Params("id")
	field := ctx.Query("field")

	var requestBody map[string]interface{}
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	id, err := strconv.Atoi(adminID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid admin ID"})
	}

	admin, err := h.strg.Admin().GetById(ctx.Context(), id)
	if admin == nil || err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin not found"})
	}

	switch field {
	case "status":
		status, ok := requestBody["status"].(string)
		if !ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status format"})
		}
		if err := h.strg.Admin().UpdateStatus(ctx.Context(), id, status); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update admin status"})
		}
		return ctx.JSON(fiber.Map{"message": "Admin status updated successfully"})

	case "limit":
		limit, ok := requestBody["limit"].(float64)
		if !ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit format"})
		}
		err := h.strg.AdminRestaurantsLimit().Update(ctx.Context(), &repo.CreateAdminRestaurantsLimit{
			AdminId:        id,
			MaxRestaurants: int(limit),
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update admin limit"})
		}
		return ctx.JSON(fiber.Map{"message": "Admin limit updated successfully"})

	default:
		var updateData repo.UpdateAdmin
		if firstName, ok := requestBody["first_name"].(string); ok {
			updateData.FirstName = firstName
		}
		if lastName, ok := requestBody["last_name"].(string); ok {
			updateData.LastName = lastName
		}
		if email, ok := requestBody["email"].(string); ok {
			updateData.Email = email
		}
		if phoneNumber, ok := requestBody["phone_number"].(string); ok {
			updateData.PhoneNumber = phoneNumber
		}
		if username, ok := requestBody["username"].(string); ok {
			updateData.Username = username
		}
		if password, ok := requestBody["password"].(string); ok && password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
			}
			updateData.PasswordHash = string(hashedPassword)
		}
		if err := h.strg.Admin().Update(ctx.Context(), id, &updateData); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update admin"})
		}
		return ctx.JSON(fiber.Map{"message": "Admin updated successfully"})
	}
}

func (h *handlerV1) GetAdminDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	admin, err := h.strg.Admin().GetById(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin not found"})
	}

	adminLogins, err := h.strg.Token().GetByAdminId(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching admin logins"})
	}

	limit, err := h.strg.AdminRestaurantsLimit().GetByAdminId(c.Context(), intID)
	if err != nil || limit == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Limit does not exist"})
	}

	restaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), admin.Id, limit.MaxRestaurants)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching admin restaurants"})
	}

	limits, err := h.strg.AdminRestaurantsLimit().GetByAdminId(c.Context(), admin.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching admin limits"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"admin":             admin,
		"admin_logins":      adminLogins,
		"admin_restaurants": restaurants,
		"admin_limits":      limits,
	})
}

func (h *handlerV1) DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	admin, err := h.strg.Admin().GetById(c.Context(), intID)
	if err != nil || admin == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin not found"})
	}

	err = h.strg.Admin().DeleteById(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Delete admin failed"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Admin deleted successfully!"})
}

func (h *handlerV1) GetProfile(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	intID, ok := adminID.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	admin, err := h.strg.Admin().GetById(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	admin.PasswordHash = ""
	return c.Status(fiber.StatusOK).JSON(admin)
}

// package v1

// import (
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
// 	"golang.org/x/crypto/bcrypt"
// )

// func (h *handlerV1) CreateAdmin(ctx *gin.Context) {
// 	var req models.CreateAdmin
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	existingAdmin, err := h.strg.Admin().GetByUsername(ctx, req.Username)
// 	if err == nil && existingAdmin != nil {
// 		ctx.JSON(http.StatusConflict, gin.H{
// 			"error": "This username already exists",
// 		})
// 		return
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error hashing password",
// 		})
// 		return
// 	}

// 	admin, err := h.strg.Admin().Create(ctx, &repo.CreateAdmin{
// 		FirstName:    req.FirstName,
// 		LastName:     req.LastName,
// 		Email:        req.Email,
// 		PhoneNumber:  req.PhoneNumber,
// 		Username:     req.Username,
// 		PasswordHash: string(hashedPassword),
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error :(",
// 		})
// 		return
// 	}

// 	NewAdmin, err := h.strg.Admin().GetByUsername(ctx, req.Username)
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error :(",
// 		})
// 		return
// 	}
// 	_, err = h.strg.AdminRestaurantsLimit().Create(ctx, &repo.CreateAdminRestaurantsLimit{
// 		AdminId:        NewAdmin.Id,
// 		MaxRestaurants: 1,
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error :(",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, models.CreateAdmin{
// 		Username: admin.Username,
// 	})
// }

// func (h *handlerV1) GetAdmins(ctx *gin.Context) {
// 	adminID := ctx.Param("id")
// 	if adminID != "" {
// 		// ID bo‘yicha bitta adminni olish
// 		id, err := strconv.Atoi(adminID)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
// 			return
// 		}

// 		admin, err := h.strg.Admin().GetById(ctx, id)
// 		if err != nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, admin)
// 		return
// 	}
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

// 	// Query parametrlarini olish
// 	status := ctx.Query("status")
// 	firstname := ctx.Query("firstname")     // firstname qidiruvi
// 	lastname := ctx.Query("lastname")       // lastname qidiruvi
// 	email := ctx.Query("email")             // email qidiruvi
// 	phonenumber := ctx.Query("phonenumber") // phonenumber qidiruvi
// 	username := ctx.Query("username")       // username qidiruvi

// 	// Sahifani 1 dan boshlash
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 10
// 	}

// 	// Adminlarni olish
// 	admins, currentPage, totalPage, err := h.strg.Admin().GetAll(
// 		ctx,
// 		status,
// 		firstname,   // firstname bo'yicha filter
// 		lastname,    // lastname bo'yicha filter
// 		email,       // email bo'yicha filter
// 		phonenumber, // phonenumber bo'yicha filter
// 		username,    // username bo'yicha filter
// 		page,
// 		limit,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admins"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":       currentPage,
// 		"total_page": totalPage,
// 		"admins":     admins,
// 	})
// }

// func (h *handlerV1) UpdateAdmin(ctx *gin.Context) {
// 	adminID := ctx.Param("id")
// 	field := ctx.Query("field") // "status", "limit" yoki umumiy yangilash

// 	var requestBody map[string]interface{}
// 	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	id, err := strconv.Atoi(adminID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
// 		return
// 	}

// 	admin, err := h.strg.Admin().GetById(ctx, id)
// 	if admin == nil || err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
// 		return
// 	}

// 	switch field {
// 	case "status":
// 		status, ok := requestBody["status"].(string)
// 		if !ok {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status format"})
// 			return
// 		}
// 		if err := h.strg.Admin().UpdateStatus(ctx, id, status); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin status"})
// 			return
// 		}
// 		if err := h.strg.Token().DeleteByAdminId(ctx, admin.Id); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke admin tokens"})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Admin status updated successfully"})
// 		return

// 	case "limit":
// 		limit, ok := requestBody["limit"].(float64)
// 		if !ok {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
// 			return
// 		}
// 		err := h.strg.AdminRestaurantsLimit().Update(ctx, &repo.CreateAdminRestaurantsLimit{
// 			AdminId:        id,
// 			MaxRestaurants: int(limit),
// 		})
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin limit"})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Admin limit updated successfully"})
// 		return

// 	default:
// 		// To'liq admin ma'lumotlarini yangilash
// 		var updateData repo.UpdateAdmin
// 		if firstName, ok := requestBody["first_name"].(string); ok {
// 			updateData.FirstName = firstName
// 		}
// 		if lastName, ok := requestBody["last_name"].(string); ok {
// 			updateData.LastName = lastName
// 		}
// 		if email, ok := requestBody["email"].(string); ok {
// 			updateData.Email = email
// 		}
// 		if phoneNumber, ok := requestBody["phone_number"].(string); ok {
// 			updateData.PhoneNumber = phoneNumber
// 		}
// 		if username, ok := requestBody["username"].(string); ok {
// 			updateData.Username = username
// 		}
// 		if password, ok := requestBody["password"].(string); ok && password != "" {
// 			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 			if err != nil {
// 				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 				return
// 			}
// 			updateData.PasswordHash = string(hashedPassword)
// 		}

// 		if err := h.strg.Admin().Update(ctx, id, &updateData); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin"})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
// 	}
// }

// func (h *handlerV1) GetAdminDetails(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	intID, err := strconv.Atoi(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}

// 	admin, err := h.strg.Admin().GetById(ctx, intID)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
// 		return
// 	}
// 	adminLogins, err := h.strg.Token().GetByAdminId(ctx, intID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching admin logins"})
// 		return
// 	}
// 	limit, err := h.strg.AdminRestaurantsLimit().GetByAdminId(ctx, intID)
// 	if err != nil || limit == nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Limit is not exists",
// 		})
// 		return
// 	}
// 	restaurants, err := h.strg.Restaurants().GetByOwnerId(ctx, admin.Id, limit.MaxRestaurants)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching admin restourants"})
// 		return
// 	}
// 	limits, err := h.strg.AdminRestaurantsLimit().GetByAdminId(ctx, admin.Id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching admin limits"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"admin":             admin,
// 		"admin_logins":      adminLogins,
// 		"admin_restourants": restaurants,
// 		"admin_limits":      limits,
// 	})
// }

// func (h *handlerV1) DeleteAdmin(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	intID, err := strconv.Atoi(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}

// 	admin, err := h.strg.Admin().GetById(ctx, intID)
// 	if err != nil || admin == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
// 		return
// 	}
// 	err = h.strg.Admin().DeleteById(ctx, intID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Delete admin failed"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Admin deleted succsessfully!"})
// }

// func (h *handlerV1) GetProfile(ctx *gin.Context) {
// 	admin_id, _ := ctx.Get("admin_id")
// 	admin, err := h.strg.Admin().GetById(ctx, admin_id.(int))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}
// 	admin.PasswordHash = ""
// 	ctx.JSON(http.StatusOK, admin)
// }
