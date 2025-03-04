package repo

import (
	"context"
	"time"
)

type TokenStorageI interface {
	Create(ctx context.Context, req *Token) (*Token, error)
	GetByAdminId(ctx context.Context, id int) ([]Token, error)
	Delete(ctx context.Context, id int) error
	DeleteByAdminId(ctx context.Context, id int) error
}

// Tokens jadvali
type Token struct {
	Id       uint      `gorm:"primaryKey"`
	AdminId  uint      `gorm:"not null"`
	Token    string    `gorm:"size:255;not null"`
	AuthTime time.Time `gorm:"autoCreateTime"`
	Admin    *Admin    `gorm:"foreignKey:AdminId;constraint:OnDelete:CASCADE"`
}
