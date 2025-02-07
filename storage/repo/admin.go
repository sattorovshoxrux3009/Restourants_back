package repo

import (
	"context"
	"time"
)

type AdminStorageI interface {
	Create(ctx context.Context, req *CreateAdmin) (*CreateAdmin, error)
	GetAll(ctx context.Context) ([]Admin, error)
	GetByUsername(ctx context.Context, username string) (*Admin, error)
	UpdatePassword(ctx context.Context, username, password string) error
	DeleteById(ctx context.Context, id int) error
	DeleteByUsername(ctx context.Context, username string) error
}
type Admin struct {
	Id           int
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	Status       string
}
type CreateAdmin struct {
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	Username     string
	PasswordHash string
}
