package v1

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

func saveImage(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	dst := filepath.Join("uploads", "restourants", newFileName)

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return "", err
	}

	if err := ctx.SaveFile(file, dst); err != nil {
		return "", err
	}

	imageURL := "/uploads/restourants/" + newFileName
	return imageURL, nil
}

func (h *handlerV1) CreateRestaurant(c *fiber.Ctx) error {
	var req models.CreateRestourants

	req.Name = c.FormValue("name")
	req.Address = c.FormValue("address")
	req.Latitude, _ = strconv.ParseFloat(c.FormValue("latitude"), 64)
	req.Longitude, _ = strconv.ParseFloat(c.FormValue("longitude"), 64)
	req.PhoneNumber = c.FormValue("phone_number")
	req.Email = c.FormValue("email")
	req.Capacity, _ = strconv.Atoi(c.FormValue("capacity"))
	req.OwnerID, _ = strconv.Atoi(c.FormValue("owner_id"))
	req.OpeningHours = c.FormValue("opening_hours")
	req.Description = c.FormValue("description")
	req.AlcoholPermission, _ = strconv.ParseBool(c.FormValue("alcohol_permission"))

	admin, err := h.strg.Admin().GetById(c.Context(), req.OwnerID)
	if err != nil || admin == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Owner ID does not exist",
		})
	}

	limit, err := h.strg.AdminRestaurantsLimit().GetByAdminId(c.Context(), req.OwnerID)
	if err != nil || limit == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Limit does not exist",
		})
	}

	ownerRestaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), req.OwnerID, limit.MaxRestaurants)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error getting owner restaurants for limit",
		})
	}

	if len(ownerRestaurants) >= limit.MaxRestaurants {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You have reached the maximum number of restaurants allowed.",
		})
	}

	file, err := c.FormFile("image")
	if err == nil && file != nil {
		imageURL, err := saveImage(c, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error saving image",
			})
		}
		req.Image = imageURL
	}

	_, err = h.strg.Restaurants().Create(c.Context(), &repo.CreateRestaurant{
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating restaurant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Creating restaurant successfully",
	})
}

func (h *handlerV1) GetRestourants(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	if restaurantID != "" {
		// ID bo‘yicha bitta restoranni olish
		id, err := strconv.Atoi(restaurantID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
		}

		restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
		if err != nil || restaurant == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
		}
		if restaurant.Status != "active" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
		}
		restaurant.Status = ""
		return c.Status(fiber.StatusOK).JSON(restaurant)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	name := c.Query("name")
	address := c.Query("address")
	capacity := c.Query("capacity")
	alcoholPermission := c.Query("adlcohol_permission")

	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetAll(
		c.Context(),
		name,
		address,
		capacity,
		alcoholPermission,
		page,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get restaurants"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":        currentPage,
		"total_page":  totalPage,
		"restaurants": restaurants,
	})
}

func (h *handlerV1) GetSRestourants(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	if restaurantID != "" {
		id, err := strconv.Atoi(restaurantID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
		}

		restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
		if err != nil || restaurant == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
		}
		return c.JSON(restaurant)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	status := c.Query("status")
	phonenumber := c.Query("phonenumber")
	email := c.Query("email")
	ownerid := c.Query("ownerid")
	name := c.Query("name")
	address := c.Query("address")
	capacity := c.Query("capacity")
	alcoholPermission := c.Query("alcoholpermission")

	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetSall(
		c.Context(),
		status,
		phonenumber,
		email,
		ownerid,
		name,
		address,
		capacity,
		alcoholPermission,
		page,
		limit,
	)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get restaurants"})
	}

	return c.JSON(fiber.Map{
		"page":        currentPage,
		"total_page":  totalPage,
		"restaurants": restaurants,
	})
}

func (h *handlerV1) UpdateRestaurantStatus(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	if restaurantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Restaurant ID is required"})
	}

	id, err := strconv.Atoi(restaurantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
	}

	var requestBody models.UpdateRestaurantStatus
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
	}

	err = h.strg.Restaurants().UpdateStatus(c.Context(), id, requestBody.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update restaurant status"})
	}

	return c.JSON(fiber.Map{"message": "Restaurant status updated successfully"})
}

func (h *handlerV1) UpdateRestaurant(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	if restaurantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Restaurant ID is required"})
	}

	id, err := strconv.Atoi(restaurantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
	}

	// Restoranni bazadan olish
	restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
	}

	// Form-data dan ma'lumotlarni olish
	var req models.UpdateRestaurants
	req.Name = c.FormValue("name")
	req.Address = c.FormValue("address")
	req.Latitude, _ = strconv.ParseFloat(c.FormValue("latitude"), 64)
	req.Longitude, _ = strconv.ParseFloat(c.FormValue("longitude"), 64)
	req.PhoneNumber = c.FormValue("phone_number")
	req.Email = c.FormValue("email")
	req.Capacity, _ = strconv.Atoi(c.FormValue("capacity"))
	req.OwnerID, _ = strconv.Atoi(c.FormValue("owner_id"))
	req.OpeningHours = c.FormValue("opening_hours")
	req.Description = c.FormValue("description")
	req.AlcoholPermission, _ = strconv.ParseBool(c.FormValue("alcohol_permission"))

	// Yangi rasm bor yoki yo‘qligini tekshirish
	file, err := c.FormFile("image")
	if err == nil {
		// Eski rasmni o‘chirish
		if restaurant.ImageURL != "" {
			oldImagePath := filepath.Join("uploads", "restaurants", filepath.Base(restaurant.ImageURL))
			_ = os.Remove(oldImagePath)
		}

		// Yangi rasmni saqlash
		imageURL, err := saveImage(c, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving new image"})
		}
		req.ImageURL = imageURL
	} else {
		req.ImageURL = restaurant.ImageURL // Agar yangi rasm kelmasa, eski rasmni saqlaymiz
	}

	err = h.strg.Restaurants().Update(c.Context(), id, &repo.UpdateRestaurant{
		Name:              req.Name,
		Address:           req.Address,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		PhoneNumber:       req.PhoneNumber,
		Email:             req.Email,
		Capacity:          req.Capacity,
		OwnerID:           req.OwnerID,
		OpeningHours:      req.OpeningHours,
		ImageURL:          req.ImageURL,
		Description:       req.Description,
		AlcoholPermission: req.AlcoholPermission,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update restaurant"})
	}

	return c.JSON(fiber.Map{"message": "Restaurant updated successfully"})
}

func (h *handlerV1) GetSRestaurantDetails(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	id, err := strconv.Atoi(restaurantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
	}

	restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
	}

	restaurantMenu, err := h.strg.Menu().GetByRestourantId(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get restaurant menu"})
	}

	return c.JSON(fiber.Map{
		"restaurant":      restaurant,
		"restaurant_menu": restaurantMenu,
	})
}

func (h *handlerV1) DeleteRestaurant(c *fiber.Ctx) error {
	restaurantID := c.Params("id")
	if restaurantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Restaurant ID is required"})
	}

	id, err := strconv.Atoi(restaurantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid restaurant ID"})
	}

	restaurant, err := h.strg.Restaurants().GetById(c.Context(), id)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
	}

	if restaurant.ImageURL != "" {
		oldImagePath := filepath.Join("uploads", "restaurants", filepath.Base(restaurant.ImageURL))
		_ = os.Remove(oldImagePath)
	}

	err = h.strg.Restaurants().Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete restaurant"})
	}

	return c.JSON(fiber.Map{"message": "Restaurant deleted successfully"})
}


// package v1

// import (
// 	"fmt"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
// )

// func saveImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
// 	fileExtension := filepath.Ext(file.Filename)
// 	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
// 	dst := filepath.Join("uploads", "restourants", newFileName)
// 	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := ctx.SaveUploadedFile(file, dst); err != nil {
// 		return "", err
// 	}
// 	imageURL := "/uploads/restourants/" + newFileName
// 	return imageURL, nil
// }

// func (h *handlerV1) CreateRestaurant(ctx *gin.Context) {
// 	var req models.CreateRestourants

// 	req.Name = ctx.PostForm("name")
// 	req.Address = ctx.PostForm("address")
// 	req.Latitude, _ = strconv.ParseFloat(ctx.PostForm("latitude"), 64)
// 	req.Longitude, _ = strconv.ParseFloat(ctx.PostForm("longitude"), 64)
// 	req.PhoneNumber = ctx.PostForm("phone_number")
// 	req.Email = ctx.PostForm("email")
// 	req.Capacity, _ = strconv.Atoi(ctx.PostForm("capacity"))
// 	req.OwnerID, _ = strconv.Atoi(ctx.PostForm("owner_id"))
// 	req.OpeningHours = ctx.PostForm("opening_hours")
// 	req.Description = ctx.PostForm("description")
// 	req.AlcoholPermission, _ = strconv.ParseBool(ctx.PostForm("alcohol_permission"))

// 	admin, err := h.strg.Admin().GetById(ctx, req.OwnerID)
// 	if err != nil || admin == nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Owner id is not exists",
// 		})
// 		return
// 	}
// 	limit, err := h.strg.AdminRestaurantsLimit().GetByAdminId(ctx, req.OwnerID)
// 	if err != nil || limit == nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Limit is not exists",
// 		})
// 		return
// 	}
// 	ownerRestaurants, err := h.strg.Restaurants().GetByOwnerId(ctx, req.OwnerID, limit.MaxRestaurants)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Error getting owner restourants for limit",
// 		})
// 		return
// 	}
// 	if len(ownerRestaurants) >= limit.MaxRestaurants {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "You have reached the maximum number of restaurants allowed.",
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
// 		imageURL, err := saveImage(ctx, file)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"error": "Error saving image",
// 			})
// 			return
// 		}
// 		req.Image = imageURL
// 	}

// 	_, err = h.strg.Restaurants().Create(ctx, &repo.CreateRestaurant{
// 		Name:              req.Name,
// 		Address:           req.Address,
// 		Latitude:          req.Latitude,
// 		Longitude:         req.Longitude,
// 		PhoneNumber:       req.PhoneNumber,
// 		Email:             req.Email,
// 		Capacity:          req.Capacity,
// 		OwnerID:           req.OwnerID,
// 		OpeningHours:      req.OpeningHours,
// 		ImageURL:          req.Image,
// 		Description:       req.Description,
// 		AlcoholPermission: req.AlcoholPermission,
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error creating restaurant",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"message": "Creating restaurant sucsessfully",
// 	})
// }

// func (h *handlerV1) GetRestourants(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	if restaurantID != "" {
// 		// ID bo‘yicha bitta adminni olish
// 		id, err := strconv.Atoi(restaurantID)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 			return
// 		}

// 		restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 		if err != nil || restaurant == nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 			return
// 		}
// 		if restaurant.Status != "active" {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 			return
// 		}
// 		restaurant.Status = ""
// 		ctx.JSON(http.StatusOK, restaurant)
// 		return
// 	}
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 10
// 	}
// 	name := ctx.Query("name")
// 	address := ctx.Query("address")
// 	capacity := ctx.Query("capacity")
// 	adlcohol_permission := ctx.Query("adlcohol_permission")

// 	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetAll(
// 		ctx,
// 		name,
// 		address,
// 		capacity,
// 		adlcohol_permission,
// 		page,
// 		limit,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurants"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":        currentPage,
// 		"total_page":  totalPage,
// 		"restaurants": restaurants,
// 	})
// }

// func (h *handlerV1) GetSRestourants(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	if restaurantID != "" {
// 		id, err := strconv.Atoi(restaurantID)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 			return
// 		}

// 		restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 		if err != nil || restaurant == nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, restaurant)
// 		return
// 	}
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 10
// 	}
// 	status := ctx.Query("status")
// 	phonenumber := ctx.Query("phonenumber")
// 	email := ctx.Query("email")
// 	ownerid := ctx.Query("ownerid")
// 	name := ctx.Query("name")
// 	address := ctx.Query("address")
// 	capacity := ctx.Query("capacity")
// 	adlcoholpermission := ctx.Query("alcoholpermission")

// 	restaurants, currentPage, totalPage, err := h.strg.Restaurants().GetSall(
// 		ctx,
// 		status,
// 		phonenumber,
// 		email,
// 		ownerid,
// 		name,
// 		address,
// 		capacity,
// 		adlcoholpermission,
// 		page,
// 		limit,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurants"})
// 		fmt.Println(err)
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":        currentPage,
// 		"total_page":  totalPage,
// 		"restaurants": restaurants,
// 	})
// }

// func (h *handlerV1) UpdateRestaurantStatus(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	if restaurantID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
// 		return
// 	}
// 	id, err := strconv.Atoi(restaurantID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 		return
// 	}
// 	var requestBody models.UpdateRestaurantStatus
// 	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}
// 	restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 	if err != nil || restaurant == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 		return
// 	}
// 	err = h.strg.Restaurants().UpdateStatus(ctx, id, requestBody.Status)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant status"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Restaurant status updated successfully"})
// }

// func (h *handlerV1) UpdateRestaurant(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	if restaurantID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
// 		return
// 	}

// 	id, err := strconv.Atoi(restaurantID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 		return
// 	}

// 	// Restoranni bazadan olish
// 	restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 	if err != nil || restaurant == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 		return
// 	}

// 	// Form-data dan ma'lumotlarni olish
// 	var req models.UpdateRestaurants
// 	req.Name = ctx.PostForm("name")
// 	req.Address = ctx.PostForm("address")
// 	req.Latitude, _ = strconv.ParseFloat(ctx.PostForm("latitude"), 64)
// 	req.Longitude, _ = strconv.ParseFloat(ctx.PostForm("longitude"), 64)
// 	req.PhoneNumber = ctx.PostForm("phone_number")
// 	req.Email = ctx.PostForm("email")
// 	req.Capacity, _ = strconv.Atoi(ctx.PostForm("capacity"))
// 	req.OwnerID, _ = strconv.Atoi(ctx.PostForm("owner_id"))
// 	req.OpeningHours = ctx.PostForm("opening_hours")
// 	req.Description = ctx.PostForm("description")
// 	req.AlcoholPermission, _ = strconv.ParseBool(ctx.PostForm("alcohol_permission"))

// 	// Yangi rasm bor yoki yo‘qligini tekshirish
// 	file, _ := ctx.FormFile("image")
// 	if file != nil {
// 		// Eskisini o‘chirish
// 		if restaurant.ImageURL != "" {
// 			oldImagePath := filepath.Join("uploads", "restourants", filepath.Base(restaurant.ImageURL))
// 			_ = os.Remove(oldImagePath) // Eskisini o‘chirish
// 		}

// 		// Yangi rasmni saqlash
// 		imageURL, err := saveImage(ctx, file)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving new image"})
// 			return
// 		}
// 		req.ImageURL = imageURL
// 	} else {
// 		req.ImageURL = restaurant.ImageURL // Agar yangi rasm kelmasa, eski rasmni saqlaymiz
// 	}
// 	err = h.strg.Restaurants().Update(ctx, id, &repo.UpdateRestaurant{
// 		Name:              req.Name,
// 		Address:           req.Address,
// 		Latitude:          req.Latitude,
// 		Longitude:         req.Longitude,
// 		PhoneNumber:       req.PhoneNumber,
// 		Email:             req.Email,
// 		Capacity:          req.Capacity,
// 		OwnerID:           req.OwnerID,
// 		OpeningHours:      req.OpeningHours,
// 		ImageURL:          req.ImageURL,
// 		Description:       req.Description,
// 		AlcoholPermission: req.AlcoholPermission,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
// }

// func (h *handlerV1) GetSRestaurantDetails(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	id, err := strconv.Atoi(restaurantID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 		return
// 	}

// 	restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 	if err != nil || restaurant == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 		return
// 	}

// 	restaurant_menu, err := h.strg.Menu().GetByRestourantId(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurant menu"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"restaurant":      restaurant,
// 		"restaurant_menu": restaurant_menu,
// 	})
// }

// func (h *handlerV1) DeleteRastaurant(ctx *gin.Context) {
// 	restaurantID := ctx.Param("id")
// 	if restaurantID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
// 		return
// 	}

// 	id, err := strconv.Atoi(restaurantID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
// 		return
// 	}

// 	restaurant, err := h.strg.Restaurants().GetById(ctx, id)
// 	if err != nil || restaurant == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
// 		return
// 	}

// 	if restaurant.ImageURL != "" {
// 		oldImagePath := filepath.Join("uploads", "restourants", filepath.Base(restaurant.ImageURL))
// 		_ = os.Remove(oldImagePath)
// 	}

// 	err = h.strg.Restaurants().Delete(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restaurant"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
// }
