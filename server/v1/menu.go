package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

func saveMenuImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	dst := filepath.Join("uploads", "menu", newFileName)
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return "", err
	}
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	imageURL := "/uploads/menu/" + newFileName
	return imageURL, nil
}

func (h *handlerV1) CreateMenu(ctx *gin.Context) {
	var req models.CreateMenu

	req.RestaurantId, _ = strconv.Atoi(ctx.PostForm("restaurant_id"))
	req.Name = ctx.PostForm("name")
	req.Description = ctx.PostForm("description")
	req.Price, _ = strconv.ParseFloat(ctx.PostForm("price"), 64)
	req.Category = ctx.PostForm("category")

	restaurant, err := h.strg.Restaurants().GetById(ctx, req.RestaurantId)
	if restaurant == nil || err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Restaurant is not exists in this Id",
		})
		return
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Form-data error",
		})
		return
	}

	file, _ := ctx.FormFile("image")
	if file != nil {
		imageURL, err := saveMenuImage(ctx, file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error saving image",
			})
			return
		}
		req.Image = imageURL
	}

	menu, err := h.strg.Menu().Create(ctx, &repo.CreateMenu{
		RestaurantId: req.RestaurantId,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Category:     req.Category,
		ImageURL:     req.Image,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error :(",
		})
		return
	}
	ctx.JSON(http.StatusCreated, models.CreateMenu{
		RestaurantId: menu.RestaurantId,
		Name:         menu.Name,
		Description:  menu.Description,
		Price:        menu.Price,
		Category:     menu.Category,
		Image:        menu.ImageURL,
	})
}

func (h *handlerV1) GetMenu(ctx *gin.Context) {
	menuId := ctx.Param("id")
	if menuId != "" {
		// ID bo‘yicha bitta menuni olish
		id, err := strconv.Atoi(menuId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
			return
		}

		menu, err := h.strg.Menu().GetById(ctx, id)
		if err != nil || menu == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}
		menu.CreatedAt = time.Time{}
		menu.UpdatedAt = time.Time{}
		ctx.JSON(http.StatusOK, menu)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	name := ctx.Query("name")
	category := ctx.Query("category")

	menu, currentPage, totalPage, err := h.strg.Menu().GetAll(
		ctx,
		name,
		category,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menus"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page":       currentPage,
		"total_page": totalPage,
		"menu":       menu,
	})

}

func (h *handlerV1) GetSMenu(ctx *gin.Context) {
	menuId := ctx.Param("id")
	if menuId != "" {
		// ID bo‘yicha bitta menuni olish
		id, err := strconv.Atoi(menuId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
			return
		}

		menu, err := h.strg.Menu().GetById(ctx, id)
		if err != nil || menu == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}
		ctx.JSON(http.StatusOK, menu)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	name := ctx.Query("name")
	category := ctx.Query("category")

	menu, currentPage, totalPage, err := h.strg.Menu().GetSAll(
		ctx,
		name,
		category,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menus"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page":       currentPage,
		"total_page": totalPage,
		"menu":       menu,
	})

}
