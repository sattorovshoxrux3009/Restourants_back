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
