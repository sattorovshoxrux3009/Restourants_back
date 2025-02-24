package repo

import (
	"context"
	"time"
)

type MenuI interface {
	Create(ctx context.Context, req *CreateMenu) (*CreateMenu, error)
	GetAll(ctx context.Context, name, category string, page, limit int) ([]Menu, int, int, error)
	GetSAll(ctx context.Context, name, category string, restaurant_id, page, limit int) ([]MenuWithStatus, int, int, error)
	GetById(ctx context.Context, id int) (*Menu, error)
	Update(ctx context.Context, id int, req *CreateMenu) (*CreateMenu, error)
	Delete(ctx context.Context, id int) error
}

type Menu struct {
	Id           uint64
	RestaurantId int
	Name         string
	Description  string
	Price        float64
	Category     string
	ImageURL     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type MenuWithStatus struct {
	Id           int       `json:"id"`
	RestaurantId int       `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Category     string    `json:"category"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`
}

type CreateMenu struct {
	RestaurantId int
	Name         string
	Description  string
	Price        float64
	Category     string
	ImageURL     string
}
