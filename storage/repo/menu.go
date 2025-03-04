package repo

import (
	"context"
	"time"
)

type MenuI interface {
	Create(ctx context.Context, req *Menu) (*Menu, error)
	GetAll(ctx context.Context, name, category string, page, limit int) ([]Menu, int, int, error)
	GetSAll(ctx context.Context, name, category string, restaurant_id, page, limit int) ([]Menu, int, int, error)
	GetSAllByRestaurants(ctx context.Context, name, category string, restaurantIDs []uint, page, limit int) ([]Menu, int, int, error)
	GetById(ctx context.Context, id int) (*Menu, error)
	GetByRestaurantId(ctx context.Context, id int) ([]*Menu, error)
	Update(ctx context.Context, id int, req *Menu) (*Menu, error)
	Delete(ctx context.Context, id int) error
}

type Menu struct {
	Id           uint       `gorm:"primaryKey"`
	RestaurantId uint       `gorm:"not null"`
	Name         string     `gorm:"size:255;not null"`
	Description  string     `gorm:"type:text;not null"`
	Price        float64    `gorm:"type:decimal(8,2);not null"`
	ImageURL     string     `gorm:"size:255;not null"`
	Category     string     `gorm:"size:255;not null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	Restaurant   Restaurant `json:"-" gorm:"foreignKey:RestaurantId;constraint:OnDelete:CASCADE"`
	Status       string     `gorm:"column:restaurant_status"` // SQL natijasidagi restaurant_status shu yerga tushadi
}
