package repo

import (
	"context"
	"time"
)

type MenuI interface {
	Create(ctx context.Context, req *CreateMenu) (*CreateMenu, error)
	GetAll(ctx context.Context, name, category string, page, limit int) ([]Menu, int, int, error)
	GetSAll(ctx context.Context, name, category string, page, limit int) ([]Menu, int, int, error)
	GetById(ctx context.Context, id int) (*Menu, error)
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
type CreateMenu struct {
	RestaurantId int
	Name         string
	Description  string
	Price        float64
	Category     string
	ImageURL     string
}
