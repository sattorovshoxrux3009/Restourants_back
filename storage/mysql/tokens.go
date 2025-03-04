package mysql

import (
	"context"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type tokenRepo struct {
	db *gorm.DB
}

func NewTokenStorage(db *gorm.DB) repo.TokenStorageI {
	return &tokenRepo{
		db: db,
	}
}

func (t *tokenRepo) Create(ctx context.Context, req *repo.Token) (*repo.Token, error) {
	// Token qo'shish
	err := t.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (t *tokenRepo) GetByAdminId(ctx context.Context, id int) ([]repo.Token, error) {
	var tokens []repo.Token

	// admin_id bo'yicha tokenlarni olish
	err := t.db.WithContext(ctx).
		Where("admin_id = ?", id).
		Find(&tokens).Error
	if err != nil {
		return nil, err
	}

	// Agar hech qanday token bo'lmasa, bo'sh slice qaytarish
	if len(tokens) == 0 {
		return []repo.Token{}, nil
	}

	return tokens, nil
}

func (t *tokenRepo) Delete(ctx context.Context, id int) error {
	// Transaction yaratish
	tx := t.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Tokenni o'chirish
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&repo.Token{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (t *tokenRepo) DeleteByAdminId(ctx context.Context, id int) error {
	// Transaction yaratish
	tx := t.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// admin_id bo'yicha tokenni o'chirish
	err := tx.WithContext(ctx).Where("admin_id = ?", id).Delete(&repo.Token{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
