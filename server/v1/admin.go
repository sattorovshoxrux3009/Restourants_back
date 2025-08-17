package v1

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Create Admin
// @Description Create a new admin user (Super Admin only)
// @Tags Super-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param admin body models.CreateAdmin true "Admin details"
// @Success 201 {object} models.AdminResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /v1/superadmin/admin [post]
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

	admin, err := h.strg.Admin().Create(ctx.Context(), &repo.Admin{
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
	newAdmin, err := h.strg.Admin().GetByUsername(ctx.Context(), req.Username)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error :("})
	}
	_, err = h.strg.AdminRestaurantsLimit().Create(ctx.Context(), &repo.AdminRestaurantLimit{
		AdminId:        newAdmin.Id,
		MaxRestaurants: 1,
	})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error :("})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"username": admin.Username})
}

// @Summary Get All Admins
// @Description Get list of all admins (Super Admin only)
// @Tags super-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int false "Admin ID (optional)"
// @Success 200 {array} models.AdminResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /v1/superadmin/admins/{id} [get]
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

// @Summary Update Admin
// @Description Update admin details (Super Admin only)
// @Tags super-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Admin ID"
// @Param field query string false "Update field (status/limit)"
// @Param admin body object true "Admin update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/superadmin/admin/{id} [put]
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
		rest, err := h.strg.Restaurants().GetByOwnerId(ctx.Context(), id, "", 20)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get restaurants"})
		}

		for _, r := range rest {
			if err := h.strg.Restaurants().UpdateStatus(ctx.Context(), int(r.Id), status); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update restaurant status"})
			}
		}

		return ctx.JSON(fiber.Map{"message": "Admin status updated successfully"})

	case "limit":
		limit, ok := requestBody["limit"].(float64)
		if !ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit format"})
		}
		err := h.strg.AdminRestaurantsLimit().Update(ctx.Context(), &repo.AdminRestaurantLimit{
			AdminId:        uint(id),
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

	restaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", limit.MaxRestaurants)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching admin restaurants"})
	}

	limits, err := h.strg.AdminRestaurantsLimit().GetByAdminId(c.Context(), int(admin.Id))
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

// @Summary Delete Admin
// @Description Delete an admin user (Super Admin only)
// @Tags super-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Admin ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/superadmin/admin/{id} [delete]
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

// @Summary Get Admin Profile
// @Description Get current admin profile
// @Tags admin-profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.AdminResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/admin/profile [get]
func (h *handlerV1) GetProfile(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)
	admin, err := h.strg.Admin().GetById(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	admin.PasswordHash = ""
	return c.Status(fiber.StatusOK).JSON(admin)

}

func (h *handlerV1) UpdateProfile(c *fiber.Ctx) error {
	// Admin ID olish
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)

	// Request body ni parse qilish
	var req models.UpdateAdmin
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Adminni bazadan olish
	admin, err := h.strg.Admin().GetById(c.Context(), intID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get admin"})
	}

	// Agar foydalanuvchi parolni o‘zgartirmoqchi bo‘lsa, eski parolni tekshirish
	var hashedPassword string
	if len(req.NewPassword) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be 4 character minimum"})
	}
	if req.OldPassword != "" && req.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.OldPassword)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid old password"})
		}
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}
		hashedPassword = string(hashedBytes)
	} else {
		// Agar yangi parol berilmagan bo‘lsa, eski parolni ishlatish
		hashedPassword = admin.PasswordHash
	}

	// Admin profilini yangilash
	err = h.strg.Admin().Update(c.Context(), intID, &repo.UpdateAdmin{
		FirstName:    req.FirstName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		LastName:     req.LastName,
		Username:     admin.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}
