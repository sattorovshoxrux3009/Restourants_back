package v1

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

func saveMenuImage(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	dst := filepath.Join("uploads", "menu", newFileName)
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return "", err
	}
	if err := c.SaveFile(file, dst); err != nil {
		return "", err
	}
	return "/uploads/menu/" + newFileName, nil
}

func (h *handlerV1) CreateSMenu(c *fiber.Ctx) error {
	restaurantId, _ := strconv.Atoi(c.FormValue("restaurant_id"))
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	restaurant, err := h.strg.Restaurants().GetById(c.Context(), restaurantId)
	if restaurant == nil || err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Restaurant not found"})
	}

	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		imageURL, err = saveMenuImage(c, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving image"})
		}
	}

	menu, err := h.strg.Menu().Create(c.Context(), &repo.CreateMenu{
		RestaurantId: restaurantId,
		Name:         c.FormValue("name"),
		Description:  c.FormValue("description"),
		Price:        price,
		Category:     c.FormValue("category"),
		ImageURL:     imageURL,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(menu)
}

func (h *handlerV1) UpdateSMenu(c *fiber.Ctx) error {
	menuId := c.Params("id")
	if menuId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Menu ID is required"})
	}

	id, err := strconv.Atoi(menuId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID"})
	}

	// Menuni bazadan olish
	menu, err := h.strg.Menu().GetById(c.Context(), id)
	if err != nil || menu == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
	}

	// Form-data dan ma'lumotlarni olish
	var req models.CreateMenu
	req.Name = c.FormValue("name")
	req.RestaurantId, _ = strconv.Atoi(c.FormValue("restaurant_id"))
	req.Description = c.FormValue("description")
	req.Price, _ = strconv.ParseFloat(c.FormValue("price"), 64)
	req.Category = c.FormValue("category")

	file, err := c.FormFile("image")
	if err == nil {
		// Eski rasmni o'chirish
		if menu.ImageURL != "" {
			oldImagePath := filepath.Join("uploads", "menu", filepath.Base(menu.ImageURL))
			_ = os.Remove(oldImagePath)
		}
		imageURL, err := saveMenuImage(c, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving new image"})
		}
		req.Image = imageURL
	} else {
		req.Image = menu.ImageURL
	}

	newMenu, err := h.strg.Menu().Update(c.Context(), id, &repo.CreateMenu{
		Name:         req.Name,
		RestaurantId: req.RestaurantId,
		Description:  req.Description,
		ImageURL:     req.Image,
		Category:     req.Category,
		Price:        req.Price,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update menu"})
	}

	return c.JSON(fiber.Map{
		"message": "Menu updated successfully",
		"menu":    newMenu,
	})
}

func (h *handlerV1) GetMenu(c *fiber.Ctx) error {
	menuId := c.Params("id")
	if menuId != "" {
		// ID bo‘yicha bitta menuni olish
		id, err := strconv.Atoi(menuId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID"})
		}

		menu, err := h.strg.Menu().GetById(c.Context(), id)
		if err != nil || menu == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
		}

		// Vaqt maydonlarini bo‘sh qilish
		menu.CreatedAt = time.Time{}
		menu.UpdatedAt = time.Time{}

		return c.JSON(menu)
	}

	// Query parametrlarini olish
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	name := c.Query("name")
	category := c.Query("category")

	menu, currentPage, totalPage, err := h.strg.Menu().GetAll(
		c.Context(),
		name,
		category,
		page,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get menus"})
	}

	return c.JSON(fiber.Map{
		"page":       currentPage,
		"total_page": totalPage,
		"menu":       menu,
	})
}

func (h *handlerV1) GetSMenu(c *fiber.Ctx) error {
	menuId := c.Params("id")
	if menuId != "" {
		// ID bo‘yicha bitta menyuni olish
		id, err := strconv.Atoi(menuId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID"})
		}

		menu, err := h.strg.Menu().GetById(c.Context(), id)
		if err != nil || menu == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
		}
		return c.JSON(menu)
	}

	// Query parametrlarini olish
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	name := c.Query("name")
	category := c.Query("category")
	restaurantIDStr := c.Query("restaurantid")

	var restaurantID int
	if restaurantIDStr != "" {
		var err error
		restaurantID, err = strconv.Atoi(restaurantIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "restaurantid must be an integer"})
		}
	}

	menu, currentPage, totalPage, err := h.strg.Menu().GetSAll(
		c.Context(),
		name,
		category,
		restaurantID,
		page,
		limit,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get menus"})
	}

	return c.JSON(fiber.Map{
		"page":       currentPage,
		"total_page": totalPage,
		"menu":       menu,
	})
}

func (h *handlerV1) DeleteSMenu(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID"})
	}

	menu, err := h.strg.Menu().GetById(c.Context(), id)
	if err != nil || menu == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
	}

	if menu.ImageURL != "" {
		oldImagePath := filepath.Join("uploads", "menu", filepath.Base(menu.ImageURL))
		_ = os.Remove(oldImagePath)
	}

	err = h.strg.Menu().Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{"message": "Menu deleted successfully"})
}

// package v1

// import (
// 	"fmt"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
// )

// func saveMenuImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
// 	fileExtension := filepath.Ext(file.Filename)
// 	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
// 	dst := filepath.Join("uploads", "menu", newFileName)
// 	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := ctx.SaveUploadedFile(file, dst); err != nil {
// 		return "", err
// 	}
// 	imageURL := "/uploads/menu/" + newFileName
// 	return imageURL, nil
// }

// func (h *handlerV1) CreateSMenu(ctx *gin.Context) {
// 	var req models.CreateMenu

// 	req.RestaurantId, _ = strconv.Atoi(ctx.PostForm("restaurant_id"))
// 	req.Name = ctx.PostForm("name")
// 	req.Description = ctx.PostForm("description")
// 	req.Price, _ = strconv.ParseFloat(ctx.PostForm("price"), 64)
// 	req.Category = ctx.PostForm("category")

// 	restaurant, err := h.strg.Restaurants().GetById(ctx, req.RestaurantId)
// 	if restaurant == nil || err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Restaurant is not exists in this Id",
// 		})
// 		return
// 	}

// 	if err := ctx.ShouldBind(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Form-data error",
// 		})
// 		return
// 	}

// 	file, _ := ctx.FormFile("image")
// 	if file != nil {
// 		imageURL, err := saveMenuImage(ctx, file)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Error saving image",
// 			})
// 			return
// 		}
// 		req.Image = imageURL
// 	}

// 	menu, err := h.strg.Menu().Create(ctx, &repo.CreateMenu{
// 		RestaurantId: req.RestaurantId,
// 		Name:         req.Name,
// 		Description:  req.Description,
// 		Price:        req.Price,
// 		Category:     req.Category,
// 		ImageURL:     req.Image,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Internal server error :(",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, models.CreateMenu{
// 		RestaurantId: menu.RestaurantId,
// 		Name:         menu.Name,
// 		Description:  menu.Description,
// 		Price:        menu.Price,
// 		Category:     menu.Category,
// 		Image:        menu.ImageURL,
// 	})
// }

// func (h *handlerV1) UpdateSMenu(ctx *gin.Context) {
// 	menuId := ctx.Param("id")
// 	if menuId == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Menu ID is required"})
// 		return
// 	}

// 	id, err := strconv.Atoi(menuId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
// 		return
// 	}

// 	// Restoranni bazadan olish
// 	menu, err := h.strg.Menu().GetById(ctx, id)
// 	if err != nil || menu == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
// 		return
// 	}

// 	// Form-data dan ma'lumotlarni olish
// 	var req models.CreateMenu
// 	req.Name = ctx.PostForm("name")
// 	req.RestaurantId, _ = strconv.Atoi(ctx.PostForm("restaurant_id"))
// 	req.Description = ctx.PostForm("description")
// 	req.Price, _ = strconv.ParseFloat(ctx.PostForm("price"), 64)
// 	req.Category = ctx.PostForm("category")
// 	req.Image = ctx.PostForm("image")

// 	file, _ := ctx.FormFile("image")
// 	if file != nil {
// 		if menu.ImageURL != "" {
// 			oldImagePath := filepath.Join("uploads", "menu", filepath.Base(menu.ImageURL))
// 			_ = os.Remove(oldImagePath)
// 		}
// 		imageURL, err := saveMenuImage(ctx, file)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving new image"})
// 			return
// 		}
// 		req.Image = imageURL
// 	} else {
// 		req.Image = menu.ImageURL
// 	}

// 	newMenu, err := h.strg.Menu().Update(ctx, id, &repo.CreateMenu{
// 		Name:         req.Name,
// 		RestaurantId: req.RestaurantId,
// 		Description:  req.Description,
// 		ImageURL:     req.Image,
// 		Category:     req.Category,
// 		Price:        req.Price,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Menu updated successfully",
// 		"menu":    newMenu,
// 	})
// }

// func (h *handlerV1) GetMenu(ctx *gin.Context) {
// 	menuId := ctx.Param("id")
// 	if menuId != "" {
// 		// ID bo‘yicha bitta menuni olish
// 		id, err := strconv.Atoi(menuId)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
// 			return
// 		}

// 		menu, err := h.strg.Menu().GetById(ctx, id)
// 		if err != nil || menu == nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
// 			return
// 		}
// 		menu.CreatedAt = time.Time{}
// 		menu.UpdatedAt = time.Time{}
// 		ctx.JSON(http.StatusOK, menu)
// 		return
// 	}
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	name := ctx.Query("name")
// 	category := ctx.Query("category")

// 	menu, currentPage, totalPage, err := h.strg.Menu().GetAll(
// 		ctx,
// 		name,
// 		category,
// 		page,
// 		limit,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menus"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":       currentPage,
// 		"total_page": totalPage,
// 		"menu":       menu,
// 	})
// }

// func (h *handlerV1) GetSMenu(ctx *gin.Context) {
// 	menuId := ctx.Param("id")
// 	if menuId != "" {
// 		// ID bo‘yicha bitta menuni olish
// 		id, err := strconv.Atoi(menuId)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
// 			return
// 		}

// 		menu, err := h.strg.Menu().GetById(ctx, id)
// 		if err != nil || menu == nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, menu)
// 		return
// 	}
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	name := ctx.Query("name")
// 	category := ctx.Query("category")
// 	restaurantIDStr := ctx.Query("restaurantid")

// 	var restaurantID int
// 	if restaurantIDStr != "" {
// 		var err error
// 		restaurantID, err = strconv.Atoi(restaurantIDStr)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "restaurantid must be an integer"})
// 			return
// 		}
// 	}
// 	menu, currentPage, totalPage, err := h.strg.Menu().GetSAll(
// 		ctx,
// 		name,
// 		category,
// 		restaurantID,
// 		page,
// 		limit,
// 	)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menus"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":       currentPage,
// 		"total_page": totalPage,
// 		"menu":       menu,
// 	})

// }

// func (h *handlerV1) DeleteSMenu(ctx *gin.Context) {
// 	menuId := ctx.Param("id")
// 	if menuId == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Menu ID is required"})
// 		return
// 	}

// 	id, err := strconv.Atoi(menuId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
// 		return
// 	}

// 	// Restoranni bazadan olish
// 	menu, err := h.strg.Menu().GetById(ctx, id)
// 	if err != nil || menu == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
// 		return
// 	}

// 	if menu.ImageURL != "" {
// 		oldImagePath := filepath.Join("uploads", "menu", filepath.Base(menu.ImageURL))
// 		_ = os.Remove(oldImagePath)
// 	}

// 	err = h.strg.Menu().Delete(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
// }
