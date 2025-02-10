package repo

import (
	"context"
	"time"
)

type TokenStorageI interface {
	Create(ctx context.Context, req *CreateToken) (*CreateToken, error)
	GetByAdminId(ctx context.Context, id int) ([]Token, error)
	Delete(ctx context.Context, id int) error
}

type Token struct {
	Id       int
	AdminId  int
	Token    string
	AuthTime time.Time
}
type CreateToken struct {
	AdminId int
	Token   string
}
