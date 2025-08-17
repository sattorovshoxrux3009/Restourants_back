package v1

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
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
	return token.SignedString([]byte("Shoxrux1801$"))
}

func (h *handlerV1) CreateSuperAdmin(ctx *fiber.Ctx) error {
	var req models.CreateSuperAdmin
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingSuperAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx.Context(), req.Username)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if existingSuperAdmin != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This username already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	superAdmin, err := h.strg.SuperAdmin().Create(ctx.Context(), &repo.SuperAdmin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error :("})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"username":  superAdmin.Username,
		"createdAt": superAdmin.CreatedAt,
	})
}

// @Summary Login
// @Description Login for admin or super admin
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Router /v1/login [post]
func (h *handlerV1) Login(ctx *fiber.Ctx) error {
	var req models.LoginAdmin
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var role, hashedPassword, username, firstName, lastName string
	var adminId uint

	superAdmin, err := h.strg.SuperAdmin().GetByUsername(context.TODO(), req.Username)
	if err == nil && superAdmin != nil {
		role = "superadmin"
		hashedPassword = superAdmin.Password
		username = superAdmin.Username
		firstName = superAdmin.FirstName
		lastName = superAdmin.LastName
	} else {
		admin, err := h.strg.Admin().GetByUsername(context.TODO(), req.Username)
		if err == nil && admin != nil {
			if admin.Status != "active" {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your account is blocked"})
			}
			role = "admin"
			hashedPassword = admin.PasswordHash
			username = admin.Username
			adminId = admin.Id
			firstName = admin.FirstName
			lastName = admin.LastName
		} else {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	token, err := CreateJWTToken(username, role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	if role == "superadmin" {
		err = h.strg.SuperAdmin().UpdateToken(ctx.Context(), username, token)

	} else {
		tokens, err := h.strg.Token().GetByAdminId(ctx.Context(), int(adminId))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tokens"})
		}
		if len(tokens) >= 5 {
			err := h.strg.Token().Delete(ctx.Context(), int(tokens[0].Id))
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete old token"})
			}
		}
		h.strg.Token().Create(ctx.Context(), &repo.Token{AdminId: adminId, Token: token})
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update token in database"})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
		"role":  role,
		"name":  firstName + " " + lastName,
	})
}

func (h *handlerV1) GetSProfile(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	superAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx.Context(), username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return ctx.JSON(superAdmin)
}

func (h *handlerV1) UpdateSProfile(ctx *fiber.Ctx) error {
	username := ctx.Locals("sadmin_username").(string)
	var req models.UpdateSuperAdmin
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	superAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx.Context(), username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(superAdmin.Password), []byte(req.OldPassword)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid old password"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}
	err = h.strg.SuperAdmin().Update(ctx.Context(), &repo.SuperAdmin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  username,
		Password:  string(hashedPassword),
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return ctx.JSON(fiber.Map{"message": "Profile updated"})
}

// package v1

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
// 	"golang.org/x/crypto/bcrypt"
// )

// func CreateJWTToken(username, role string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"username": username,
// 		"role":     role,
// 		"exp":      time.Now().Add(time.Hour * 5).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte("Shoxrux1801$"))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }

// func (h *handlerV1) CreateSuperAdmin(ctx *gin.Context) {
// 	var req models.CreateSuperAdmin
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	existingSuperAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx, req.Username)
// 	if err == nil && existingSuperAdmin != nil {
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

// 	superAdmin, err := h.strg.SuperAdmin().Create(ctx, &repo.CreateSuperAdmin{
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		Username:  req.Username,
// 		Password:  string(hashedPassword),
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error :(",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, models.CreateSuperAdmin{
// 		Username:  superAdmin.Username,
// 		CreatedAt: superAdmin.CreatedAt,
// 	})
// }

// func (h *handlerV1) Login(ctx *gin.Context) {
// 	var req models.LoginAdmin

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid input",
// 		})
// 		return
// 	}

// 	var role string
// 	var hashedPassword string
// 	var username string
// 	var adminId int
// 	var FirstName, LastName string
// 	// 1. SuperAdmin jadvalidan tekshirish
// 	superAdmin, err := h.strg.SuperAdmin().GetByUsername(context.TODO(), req.Username)
// 	if err == nil && superAdmin != nil {
// 		role = "super_admin"
// 		hashedPassword = superAdmin.Password
// 		username = superAdmin.Username
// 		FirstName = superAdmin.FirstName
// 		LastName = superAdmin.LastName
// 	} else {
// 		// 2. Admin jadvalidan tekshirish
// 		admin, err := h.strg.Admin().GetByUsername(context.TODO(), req.Username)
// 		if err == nil && admin != nil {
// 			role = "admin"
// 			hashedPassword = admin.PasswordHash
// 			username = admin.Username
// 			adminId = admin.Id
// 			FirstName = admin.FirstName
// 			LastName = admin.LastName
// 		} else {
// 			// 3. Ikkalasida ham topilmasa, xato qaytarish
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"error": "Invalid username or password",
// 			})
// 			return
// 		}
// 	}

// 	// 4. Parolni tekshirish (faqat topilgan jadvalga tegishli foydalanuvchining parolini tekshiramiz)
// 	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid username or password",
// 		})
// 		return
// 	}
// 	// 5. JWT token yaratish
// 	token, err := CreateJWTToken(username, role)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to generate token",
// 		})
// 		return
// 	}

// 	// 6. Tokenni bazaga saqlash (qaysi jadvaldan topilganiga qarab)
// 	if role == "super_admin" {
// 		err = h.strg.SuperAdmin().UpdateToken(ctx, username, token)
// 	} else {
// 		admn, err := h.strg.Admin().GetByUsername(ctx, username)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Failed to delete old token",
// 			})
// 			return
// 		}
// 		if admn.Status == "inactive" {
// 			ctx.JSON(http.StatusForbidden, gin.H{
// 				"error": "Your account is blocked",
// 			})
// 			return
// 		}
// 		tokens, err := h.strg.Token().GetByAdminId(ctx, adminId)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Failed to fetch tokens",
// 			})
// 			return
// 		}
// 		// Agar tokenlar 5 tadan ko‘p bo‘lsa, eng eski tokenni o‘chirish
// 		if len(tokens) >= 5 {
// 			// Eng eski tokenni o‘chirish (tokens[0] eng eski bo‘ladi)
// 			err := h.strg.Token().Delete(ctx, tokens[0].Id)
// 			if err != nil {
// 				ctx.JSON(http.StatusInternalServerError, gin.H{
// 					"error": "Failed to delete old token",
// 				})
// 				return
// 			}
// 		}
// 		_, err = h.strg.Token().Create(ctx, &repo.CreateToken{
// 			AdminId: adminId,
// 			Token:   token,
// 		})
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Failed to create token",
// 			})
// 			return
// 		}
// 	}

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to update token in database",
// 		})
// 		return
// 	}

// 	// 7. Tokenni qaytarish
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"token": token,
// 		"role":  role,
// 		"name":  FirstName + " " + LastName,
// 	})
// }

// func (h *handlerV1) GetSProfile(ctx *gin.Context) {
// 	username, _ := ctx.Get("username")
// 	superAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx, username.(string))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, superAdmin)
// }

// func (h *handlerV1) UpdateSProfile(ctx *gin.Context) {
// 	username, _ := ctx.Get("sadmin_username")
// 	var req models.UpdateSuperAdmin
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid input",
// 		})
// 		return
// 	}
// 	superAdmin, err := h.strg.SuperAdmin().GetByUsername(ctx, username.(string))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error",
// 		})
// 		return
// 	}
// 	err = bcrypt.CompareHashAndPassword([]byte(superAdmin.Password), []byte(req.OldPassword))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid old password",
// 		})
// 		return
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error hashing password",
// 		})
// 		return
// 	}
// 	err = h.strg.SuperAdmin().Update(ctx, &repo.SuperAdmin{
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		Username:  username.(string),
// 		Password:  string(hashedPassword),
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Profile updated",
// 	})
// }
