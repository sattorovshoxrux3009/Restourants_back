package mysql

import (
	"context"
	"errors"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type eventPricesRepo struct {
	db *gorm.DB
}

func NewEventPricesStorage(db *gorm.DB) repo.EventPricesI {
	return &eventPricesRepo{
		db: db,
	}
}

func (e *eventPricesRepo) Create(ctx context.Context, req *repo.EventPrice) (*repo.EventPrice, error) {
	if err := e.db.WithContext(ctx).Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (e *eventPricesRepo) GetAll(ctx context.Context, restaurantID, eventType string, page, limit int) ([]*repo.EventPrice, int, int, error) {
	var eventPrices []*repo.EventPrice
	query := e.db.WithContext(ctx).Model(&repo.EventPrice{})

	if restaurantID != "" {
		query = query.Where("restaurant_id = ?", restaurantID)
	}
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	var totalCount int64
	query.Count(&totalCount)

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))
	offset := (page - 1) * limit
	query = query.Order("restaurant_id ASC").Limit(limit).Offset(offset)

	if err := query.Find(&eventPrices).Error; err != nil {
		return nil, 0, 0, err
	}

	return eventPrices, page, totalPages, nil
}

func (e *eventPricesRepo) GetByRestaurantID(ctx context.Context, restaurantID int) ([]*repo.EventPrice, error) {
	var eventPrices []*repo.EventPrice
	if err := e.db.WithContext(ctx).Where("restaurant_id = ?", restaurantID).Find(&eventPrices).Error; err != nil {
		return nil, err
	}
	return eventPrices, nil
}

func (r *eventPricesRepo) GetAllByRestaurantIDs(ctx context.Context, restaurantMap map[int]bool, eventType string, page, limit int) ([]repo.EventPrice, int, int, error) {
	var events []repo.EventPrice
	var totalRecords int64

	// 1️⃣ Restaurant ID larni slice (massiv) ga o‘giramiz
	restaurantIDs := make([]int, 0, len(restaurantMap))
	for id := range restaurantMap {
		restaurantIDs = append(restaurantIDs, id)
	}

	// 2️⃣ Jami nechta event borligini hisoblash
	query := r.db.WithContext(ctx).Model(&repo.EventPrice{}).
		Where("restaurant_id IN ?", restaurantIDs)

	// Agar eventType berilgan bo‘lsa, filtr qo‘shamiz
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	// Jami eventlarni sanash
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, 0, err
	}

	// 3️⃣ Pagination hisoblash
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	offset := (page - 1) * limit

	// 4️⃣ Eventlarni olish
	err := query.Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&events).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return events, page, totalPages, nil
}

func (e *eventPricesRepo) GetByID(ctx context.Context, id int) (*repo.EventPrice, error) {
	var event repo.EventPrice
	result := e.db.WithContext(ctx).First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &event, nil
}

func (e *eventPricesRepo) Update(ctx context.Context, event *repo.EventPrice) error {
	result := e.db.WithContext(ctx).Save(event)
	return result.Error
}

func (e *eventPricesRepo) Delete(ctx context.Context, id int) error {
	result := e.db.WithContext(ctx).Delete(&repo.EventPrice{}, id)
	return result.Error
}
