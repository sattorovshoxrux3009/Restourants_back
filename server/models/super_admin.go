package models

import "time"

type SuperAdmin struct {
	Id        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
	LastLogin time.Time `json:"last_login"`
}
type CreateSuperAdmin struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
type LoginAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UpdateSuperAdmin struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
