package repo

import (
	"context"
	"time"
)

type AdminRestaurantsLimitI interface {
	Create(ctx context.Context, req *CreateAdminRestaurantsLimit) (*CreateAdminRestaurantsLimit, error)
	GetByAdminId(ctx context.Context, id int) (*AdminRestaurantsLimit, error)
	Update(ctx context.Context, req *CreateAdminRestaurantsLimit) error
}

type AdminRestaurantsLimit struct {
	Id             int
	AdminId        int
	MaxRestaurants int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type CreateAdminRestaurantsLimit struct {
	AdminId        int
	MaxRestaurants int
}
