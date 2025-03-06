package v1

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

// Create
func (h *handlerV1) CreateSEventPrices(c *fiber.Ctx) error {
	var req models.EventPrice
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.EventType != "morning" && req.EventType != "night" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event type"})
	}
	restaurant, err := h.strg.Restaurants().GetById(c.Context(), req.RestaurantId)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant does not exist"})
	}
	existingEvents, err := h.strg.EventPrices().GetByRestaurantID(c.Context(), req.RestaurantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	for _, event := range existingEvents {
		if event.EventType == req.EventType {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Event type '%s' already exists for this restaurant", req.EventType),
			})
		}
	}
	EventPrice, err := h.strg.EventPrices().Create(c.Context(), &repo.EventPrice{
		RestaurantId:      uint(req.RestaurantId),
		EventType:         req.EventType,
		TablePrice:        req.TablePrice,
		WaiterPrice:       req.WaiterPrice,
		MaxGuests:         req.MaxGuests,
		TableSeats:        req.TableSeats,
		MaxWaiters:        req.MaxWaiters,
		AlcoholPermission: req.AlcoholPermission,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"event_price": EventPrice})
}

func (h *handlerV1) CreateAEventPrices(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)

	var req models.EventPrice
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.EventType != "morning" && req.EventType != "night" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event type"})
	}
	restaurant, err := h.strg.Restaurants().GetById(c.Context(), req.RestaurantId)
	if err != nil || restaurant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant does not exist"})
	}

	adminRestaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	var access bool
	for _, r := range adminRestaurants {
		if r.Id == uint(req.RestaurantId) {
			access = true
			break
		}
	}
	if !access {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	existingEvents, err := h.strg.EventPrices().GetByRestaurantID(c.Context(), req.RestaurantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	for _, event := range existingEvents {
		if event.EventType == req.EventType {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Event type '%s' already exists for this restaurant", req.EventType),
			})
		}
	}
	EventPrice, err := h.strg.EventPrices().Create(c.Context(), &repo.EventPrice{
		RestaurantId:      uint(req.RestaurantId),
		EventType:         req.EventType,
		TablePrice:        req.TablePrice,
		WaiterPrice:       req.WaiterPrice,
		MaxGuests:         req.MaxGuests,
		TableSeats:        req.TableSeats,
		MaxWaiters:        req.MaxWaiters,
		AlcoholPermission: req.AlcoholPermission,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"event_price": EventPrice})
}

// Get
func (h *handlerV1) GetSEventPrices(c *fiber.Ctx) error {
	eventId := c.Params("id")
	if eventId != "" {
		id, err := strconv.Atoi(eventId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
		}
		event, err := h.strg.EventPrices().GetByID(c.Context(), id)
		if err != nil || event == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}
		return c.Status(fiber.StatusOK).JSON(event)
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	restaurantID := c.Query("restaurantid")
	eventType := c.Query("eventtype")

	events, currentPage, totalPages, err := h.strg.EventPrices().GetAll(c.Context(), restaurantID, eventType, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get events"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page":       currentPage,
		"total_page": totalPages,
		"events":     events,
	})
}

func (h *handlerV1) GetAEventPrices(c *fiber.Ctx) error {
	// 1️⃣ Admin ID ni olish
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)

	// 2️⃣ ID orqali bitta eventni olish
	eventId := c.Params("id")
	if eventId != "" {
		id, err := strconv.Atoi(eventId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
		}

		event, err := h.strg.EventPrices().GetByID(c.Context(), id)
		if err != nil || event == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
		}

		// 3️⃣ Adminning restoranlarini olish
		restaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", 20)
		if err != nil || len(restaurants) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Restaurant not found"})
		}

		// 4️⃣ Admin ushbu eventga egalik qiladimi?
		var access = false
		for _, r := range restaurants {
			if r.Id == event.RestaurantId {
				access = true
				break
			}
		}

		if !access {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		// ✅ Eventni qaytarish
		return c.Status(fiber.StatusOK).JSON(event)
	}

	// 5️⃣ Adminning barcha restoranlarini olish
	restaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", 50) // 50ta limit
	if err != nil || len(restaurants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No restaurants found"})
	}

	// 6️⃣ Adminning restoranlari ro‘yxatini mapga o‘giramiz (tejalgan qidiruv uchun)
	restaurantMap := make(map[int]bool)
	for _, r := range restaurants {
		restaurantMap[int(r.Id)] = true
	}

	// 7️⃣ Query orqali kelgan ma’lumotlarni olish
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	eventType := c.Query("eventtype")

	// 9️⃣ Faqat admin egalik qiladigan restoranlarning eventlarini olish
	events, currentPage, totalPages, err := h.strg.EventPrices().GetAllByRestaurantIDs(c.Context(), restaurantMap, eventType, page, limit)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get events"})
	}

	// ✅ Eventlarni qaytarish
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"events":       events,
		"current_page": currentPage,
		"total_page":   totalPages,
	})
}

// Update
func (h *handlerV1) UpdateSEventPrices(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var req models.UpdateEventPrices
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingEvent, err := h.strg.EventPrices().GetByID(c.Context(), id)
	if err != nil || existingEvent == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event price not found"})
	}

	existingEvent.TablePrice = req.TablePrice
	existingEvent.WaiterPrice = req.WaiterPrice
	existingEvent.MaxGuests = req.MaxGuests
	existingEvent.TableSeats = req.TableSeats
	existingEvent.MaxWaiters = req.MaxWaiters
	existingEvent.AlcoholPermission = req.AlcoholPermission

	if err := h.strg.EventPrices().Update(c.Context(), existingEvent); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update event price"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Event price updated successfully",
		"event_price": existingEvent,
	})
}

func (h *handlerV1) UpdateAEventPrices(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var req models.UpdateEventPrices
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingEvent, err := h.strg.EventPrices().GetByID(c.Context(), id)
	if err != nil || existingEvent == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event price not found"})
	}

	adminRestaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	var access bool
	for _, r := range adminRestaurants {
		if r.Id == uint(existingEvent.RestaurantId) {
			access = true
			break
		}
	}
	if !access {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	existingEvent.TablePrice = req.TablePrice
	existingEvent.WaiterPrice = req.WaiterPrice
	existingEvent.MaxGuests = req.MaxGuests
	existingEvent.TableSeats = req.TableSeats
	existingEvent.MaxWaiters = req.MaxWaiters
	existingEvent.AlcoholPermission = req.AlcoholPermission

	if err := h.strg.EventPrices().Update(c.Context(), existingEvent); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update event price"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event price updated successfully",
	})
}

// Delete
func (h *handlerV1) DeleteSEventPrices(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	existingEvent, err := h.strg.EventPrices().GetByID(c.Context(), id)
	if err != nil || existingEvent == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event price not found"})
	}

	if err := h.strg.EventPrices().Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event price"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event price deleted successfully"})
}

func (h *handlerV1) DeleteAEventPrices(c *fiber.Ctx) error {
	adminID := c.Locals("admin_id")
	uintID, ok := adminID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	intID := int(uintID)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	existingEvent, err := h.strg.EventPrices().GetByID(c.Context(), id)
	if err != nil || existingEvent == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event price not found"})
	}

	adminRestaurants, err := h.strg.Restaurants().GetByOwnerId(c.Context(), intID, "", 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	var access bool
	for _, r := range adminRestaurants {
		if r.Id == uint(existingEvent.RestaurantId) {
			access = true
			break
		}
	}
	if !access {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if err := h.strg.EventPrices().Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event price"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event price deleted successfully"})
}

// package v1

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sattorovshoxrux3009/Restourants_back/server/models"
// 	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
// )

// func (h *handlerV1) CreateSEventPrices(ctx *gin.Context) {
// 	var req models.CreateEventPrices
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if req.EventType != "morning" && req.EventType != "night" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
// 		return
// 	}
// 	restaurant, err := h.strg.Restaurants().GetById(ctx, req.RestaurantId)
// 	if err != nil || restaurant == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant does not exist"})
// 		return
// 	}
// 	existingEvents, err := h.strg.EventPrices().GetByRestaurantID(ctx, req.RestaurantId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}
// 	for _, event := range existingEvents {
// 		if event.EventType == req.EventType {
// 			ctx.JSON(http.StatusConflict, gin.H{
// 				"error": fmt.Sprintf("Event type '%s' already exists for this restaurant", req.EventType),
// 			})
// 			return
// 		}
// 	}
// 	EventPrice, err := h.strg.EventPrices().Create(ctx, &repo.CreateEventPrices{
// 		RestaurantId:      req.RestaurantId,
// 		EventType:         req.EventType,
// 		TablePrice:        req.TablePrice,
// 		WaiterPrice:       req.WaiterPrice,
// 		MaxGuests:         req.MaxGuests,
// 		TableSeats:        req.TableSeats,
// 		MaxWaiters:        req.MaxWaiters,
// 		AlcoholPermission: req.AlcoholPermission,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, gin.H{"event_price": EventPrice})
// }

// func (h *handlerV1) GetSEventPrices(ctx *gin.Context) {
// 	eventId := ctx.Param("id")
// 	if eventId != "" {
// 		id, err := strconv.Atoi(eventId)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
// 			return
// 		}

// 		event, err := h.strg.EventPrices().GetByID(ctx, id)
// 		if err != nil || event == nil {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, event)
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
// 	restaurantID := ctx.Query("restaurantid")
// 	eventType := ctx.Query("eventtype")

// 	events, currentPage, totalPages, err := h.strg.EventPrices().GetAll(ctx, restaurantID, eventType, page, limit)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"page":       currentPage,
// 		"total_page": totalPages,
// 		"events":     events,
// 	})
// }

// func (h *handlerV1) UpdateSEventPrices(ctx *gin.Context) {
// 	// Parametrdan ID olish
// 	idParam := ctx.Param("id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid ID format",
// 		})
// 		return
// 	}

// 	// JSON bodyni parse qilish
// 	var req models.UpdateEventPrices
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// Bazadan event price olish
// 	existingEvent, err := h.strg.EventPrices().GetByID(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to fetch event price",
// 		})
// 		return
// 	}
// 	if existingEvent == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"error": "Event price not found",
// 		})
// 		return
// 	}

// 	// Yangi qiymatlar bilan yangilash
// 	existingEvent.TablePrice = req.TablePrice
// 	existingEvent.WaiterPrice = req.WaiterPrice
// 	existingEvent.MaxGuests = req.MaxGuests
// 	existingEvent.TableSeats = req.TableSeats
// 	existingEvent.MaxWaiters = req.MaxWaiters
// 	existingEvent.AlcoholPermission = req.AlcoholPermission

// 	// Bazada yangilash
// 	err = h.strg.EventPrices().Update(ctx, existingEvent)
// 	if err != nil {
// 		fmt.Println(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to update event price",
// 		})
// 		return
// 	}

// 	// Muvaffaqiyatli javob
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message":     "Event price updated successfully",
// 		"event_price": existingEvent,
// 	})
// }

// func (h *handlerV1) DeleteSEventPrices(ctx *gin.Context) {
// 	// Parametrdan ID olish
// 	idParam := ctx.Param("id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid ID format",
// 		})
// 		return
// 	}

// 	// Bazadan event price olish
// 	existingEvent, err := h.strg.EventPrices().GetByID(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to fetch event price",
// 		})
// 		return
// 	}
// 	if existingEvent == nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"error": "Event price not found",
// 		})
// 		return
// 	}

// 	// Transaction bilan o‘chirish
// 	err = h.strg.EventPrices().Delete(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to delete event price",
// 		})
// 		return
// 	}

// 	// Muvaffaqiyatli javob
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Event price deleted successfully",
// 	})
// }
