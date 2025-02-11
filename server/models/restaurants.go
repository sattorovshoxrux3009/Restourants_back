package models

import "time"

type Restourants struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	PhoneNumber       string    `json:"phone_number"`
	Email             string    `json:"email"`
	Capacity          int       `json:"capacity"`
	OwnerID           int       `json:"owner_id"`
	OpeningHours      string    `json:"opening_hours"`
	ImageURL          string    `json:"image_url"`
	Description       string    `json:"description"`
	AlcoholPermission bool      `json:"alcohol_permission"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateRestourants struct {
	Name              string  `json:"name"`
	Address           string  `json:"address"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	PhoneNumber       string  `json:"phone_number"`
	Email             string  `json:"email"`
	Capacity          int     `json:"capacity"`
	OwnerID           int     `json:"owner_id"`
	OpeningHours      string  `json:"opening_hours"`
	Image             string  `json:"image"`
	Description       string  `json:"description"`
	AlcoholPermission bool    `json:"alcohol_permission"`
}
