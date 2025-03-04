package repo

import (
	"context"
	"time"
)

type SuperAdminStorageI interface {
	Create(ctx context.Context, req *SuperAdmin) (*SuperAdmin, error)
	GetByUsername(ctx context.Context, username string) (*SuperAdmin, error)
	GetToken(ctx context.Context, username string) (string, error)
	Update(ctx context.Context, req *SuperAdmin) error
	UpdatePassword(ctx context.Context, username, password string) error
	UpdateToken(ctx context.Context, username, token string) error
}

type SuperAdmin struct {
	Id        uint      `gorm:"primaryKey"`
	FirstName string    `gorm:"size:255;not null"`
	LastName  string    `gorm:"size:255;not null"`
	Username  string    `gorm:"size:100;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // ✅ TIMESTAMP ishlatilmoqda
	Token     *string
	LastLogin *time.Time `gorm:"autoUpdateTime"`
}
