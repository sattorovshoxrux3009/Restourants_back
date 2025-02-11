package v1

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

func saveImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	dst := filepath.Join("uploads", "restourants", newFileName)
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return "", err
	}
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	imageURL := "/uploads/restourants/" + newFileName
	return imageURL, nil
}

func (h *handlerV1) CreateRestaurant(ctx *gin.Context) {
	var req models.CreateRestourants

	req.Name = ctx.PostForm("name")
	req.Address = ctx.PostForm("address")
	req.Latitude, _ = strconv.ParseFloat(ctx.PostForm("latitude"), 64)
	req.Longitude, _ = strconv.ParseFloat(ctx.PostForm("longitude"), 64)
	req.PhoneNumber = ctx.PostForm("phone_number")
	req.Email = ctx.PostForm("email")
	req.Capacity, _ = strconv.Atoi(ctx.PostForm("capacity"))
	req.OwnerID, _ = strconv.Atoi(ctx.PostForm("owner_id"))
	req.OpeningHours = ctx.PostForm("opening_hours")
	req.Description = ctx.PostForm("description")
	req.AlcoholPermission, _ = strconv.ParseBool(ctx.PostForm("alcohol_permission"))

	admin, err := h.strg.Admin().GetById(ctx, req.OwnerID)
	if err != nil || admin == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Owner id is not exists",
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
		imageURL, err := saveImage(ctx, file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error saving image",
			})
			return
		}
		req.Image = imageURL
	}

	restaurant, err := h.strg.Restaurants().Create(ctx, &repo.CreateRestaurant{
		Name:              req.Name,
		Address:           req.Address,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		PhoneNumber:       req.PhoneNumber,
		Email:             req.Email,
		Capacity:          req.Capacity,
		OwnerID:           req.OwnerID,
		OpeningHours:      req.OpeningHours,
		ImageURL:          req.Image,
		Description:       req.Description,
		AlcoholPermission: req.AlcoholPermission,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating restaurant",
		})
		return
	}

	ctx.JSON(http.StatusCreated, restaurant)
}

func (h *handlerV1) GetRestourants(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID != "" {
		// ID bo‘yicha bitta adminni olish
		id, err := strconv.Atoi(restaurantID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		restaurant, err := h.strg.Restaurants().GetById(ctx, id)
		if err != nil || restaurant == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		if restaurant.Status != "active" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		restaurant.Status = ""
		ctx.JSON(http.StatusOK, restaurant)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	name := ctx.Query("name")
	address := ctx.Query("address")
	capacity := ctx.Query("capacity")
	adlcohol_permission := ctx.Query("adlcohol_permission")

	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetAll(
		ctx,
		name,
		address,
		capacity,
		adlcohol_permission,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurants"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page":        currentPage,
		"total_page":  totalPage,
		"restaurants": restaurants,
	})
}

func (h *handlerV1) GetSRestourants(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID != "" {
		id, err := strconv.Atoi(restaurantID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		restaurant, err := h.strg.Restaurants().GetById(ctx, id)
		if err != nil || restaurant == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		ctx.JSON(http.StatusOK, restaurant)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	status := ctx.Query("status")
	name := ctx.Query("name")
	address := ctx.Query("address")
	capacity := ctx.Query("capacity")
	adlcohol_permission := ctx.Query("adlcohol_permission")

	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetSall(
		ctx,
		status,
		name,
		address,
		capacity,
		adlcohol_permission,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurants"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page":        currentPage,
		"total_page":  totalPage,
		"restaurants": restaurants,
	})
}
