package repo

import (
	"context"
	"time"
)

type RestaurantsI interface {
	Create(ctx context.Context, req *Restaurant) (*Restaurant, error)
	GetAll(ctx context.Context, name, address, capacity, adlcohol_permission string, page, limit int) ([]Restaurant, int, int, error)
	GetSall(ctx context.Context, status, phonenumber, email, ownerid, name, address, capacity, alcohol_permission string, page, limit int) ([]Restaurant, int, int, error)
	GetByOwnerId(ctx context.Context, ownerID int, name string, limit int) ([]Restaurant, error)
	GetById(ctx context.Context, id int) (*Restaurant, error)
	Update(ctx context.Context, id int, req *Restaurant) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
	SearchByName(ctx context.Context, nameQuery string, page int, limit int) ([]RestaurantShort, int, error)
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

// Restaurants jadvali
type Restaurant struct {
	Id                uint      `gorm:"primaryKey"`
	Name              string    `gorm:"size:255;not null"`
	Address           string    `gorm:"size:255;not null"`
	Latitude          float64   `gorm:"type:decimal(10,7);not null"`
	Longitude         float64   `gorm:"type:decimal(10,7);not null"`
	PhoneNumber       string    `gorm:"size:20;not null"`
	Email             string    `gorm:"size:255;not null"`
	Capacity          int       `gorm:"not null"`
	OwnerId           uint      `gorm:"not null"`
	OpeningHours      string    `gorm:"size:255;not null"`
	ImageURL          string    `gorm:"size:255;not null"`
	Description       string    `gorm:"type:text;not null"`
	AlcoholPermission bool      `gorm:"not null"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	Status            string    `gorm:"type:enum('active','inactive');default:'active'"`
	Owner             Admin     `json:"-" gorm:"foreignKey:OwnerId;constraint:OnDelete:CASCADE"`
}
type RestaurantShort struct {
	Id       uint
	Name     string
	ImageURL string
	Type     string
}
