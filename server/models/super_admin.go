package models

import "time"

type SuperAdmin struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
	LastLogin time.Time `json:"last_login"`
}
type CreateSuperAdmin struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
type LoginAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
