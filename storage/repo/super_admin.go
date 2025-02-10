package repo

import (
	"context"
	"time"
)

type SuperAdminStorageI interface {
	Create(ctx context.Context, req *CreateSuperAdmin) (*CreateSuperAdmin, error)
	GetByUsername(ctx context.Context, username string) (*SuperAdmin, error)
	GetToken(ctx context.Context, username string) (string, error)
	UpdatePassword(ctx context.Context, username, password string) error
	UpdateToken(ctx context.Context, username, token string) error
}
type SuperAdmin struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
	Token     string
	LastLogin time.Time
}
type CreateSuperAdmin struct {
	Username  string
	Password  string
	CreatedAt time.Time
}
