package models

import "time"

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Message string `json:"message" example:"Bad Request"`
}

// @Description Success response structure
type SuccessResponse struct {
	Message string `json:"message" example:"Success"`
	Name    string `json:"name,omitempty" example:"Admin"`
	Role    string `json:"role,omitempty" example:"Super-admin"`
	Token   string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RestaurantResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Image       string    `json:"image"`
	AdminID     int       `json:"admin_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type MenuResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	RestaurantID int       `json:"restaurant_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type AdminResponse struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
}
