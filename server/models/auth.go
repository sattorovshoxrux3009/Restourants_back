package models

// LoginRequest defines the request body for login
// @Description Login request structure
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password123"`
}

// LoginResponse defines the response for successful login
// @Description Login response structure
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
