package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type superAdminRepo struct {
	db *gorm.DB
}

func NewSuperAdminStorage(db *gorm.DB) repo.SuperAdminStorageI {
	return &superAdminRepo{
		db: db,
	}
}

func (s *superAdminRepo) Create(ctx context.Context, req *repo.SuperAdmin) (*repo.SuperAdmin, error) {
	if err := s.db.Create(req).Error; err != nil {
		return nil, err
	}

	return req, nil
}

func (s *superAdminRepo) GetByUsername(ctx context.Context, username string) (*repo.SuperAdmin, error) {
	var admin repo.SuperAdmin
	result := s.db.WithContext(ctx).Where("username = ?", username).First(&admin)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Foydalanuvchi yo‘q bo‘lsa, shunchaki nil qaytarish
		}
		fmt.Println("DB QUERY ERROR:", result.Error) // ✅ Faqat haqiqiy xatolarni chiqarish
		return nil, result.Error
	}

	return &admin, nil
}

func (s *superAdminRepo) GetToken(ctx context.Context, username string) (string, error) {
	var token string

	err := s.db.Model(&repo.SuperAdmin{}).
		Select("token").
		Where("username = ?", username).
		Scan(&token).Error

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *superAdminRepo) Update(ctx context.Context, req *repo.SuperAdmin) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&repo.SuperAdmin{}).
			Where("username = ?", req.Username).
			Updates(map[string]interface{}{
				"first_name": req.FirstName,
				"last_name":  req.LastName,
				"password":   req.Password,
			})

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (s *superAdminRepo) UpdatePassword(ctx context.Context, username, password string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&repo.SuperAdmin{}).
			Where("username = ?", username).
			Update("password", password)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (s *superAdminRepo) UpdateToken(ctx context.Context, username, token string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&repo.SuperAdmin{}).
			Where("username = ?", username).
			Updates(map[string]interface{}{
				"token":      token,
				"last_login": time.Now(),
			})

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}
