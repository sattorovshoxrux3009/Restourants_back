package repo

import (
	"context"
	"time"
)

type RestourantsI interface {
	Create(ctx context.Context, req *CreateRestourant) (*CreateRestourant, error)
	GetByOwnerId(ctx context.Context, id int) ([]Restaurant, error)
	GetById(ctx context.Context, id int) (*Restaurant, error)
	Update(ctx context.Context, id int, req *UpdateRestourant) error
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
}
type CreateRestourant struct {
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
type UpdateRestourant struct {
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
