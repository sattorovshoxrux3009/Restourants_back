package repo

import (
	"context"
	"time"
)

type RestaurantsI interface {
	Create(ctx context.Context, req *CreateRestaurant) (*CreateRestaurant, error)
	GetAll(ctx context.Context, name, address, capacity, adlcohol_permission string, page, limit int) ([]Restaurant, int, int, error)
	GetSall(ctx context.Context, status, phonenumber, email, ownerid, name, address, capacity, alcohol_permission string, page, limit int) ([]Restaurant, int, int, error)
	GetByOwnerId(ctx context.Context, id, limit int) ([]Restaurant, error)
	GetById(ctx context.Context, id int) (*Restaurant, error)
	Update(ctx context.Context, id int, req *UpdateRestaurant) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}

type Restaurant struct {
	Id                int
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	PhoneNumber       string
	Email             string
	Capacity          int
	OwnerID           int
	OpeningHours      string
	ImageURL          string
	Description       string
	AlcoholPermission bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Status            string
}
type CreateRestaurant struct {
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	PhoneNumber       string
	Email             string
	Capacity          int
	OwnerID           int
	OpeningHours      string
	ImageURL          string
	Description       string
	AlcoholPermission bool
}
type UpdateRestaurant struct {
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	PhoneNumber       string
	Email             string
	Capacity          int
	OwnerID           int
	OpeningHours      string
	ImageURL          string
	Description       string
	AlcoholPermission bool
}
