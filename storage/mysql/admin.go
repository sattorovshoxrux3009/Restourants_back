package mysql

import (
	"context"
	"math"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type adminRepo struct {
	db *gorm.DB
}

func NewAdminStorage(db *gorm.DB) repo.AdminStorageI {
	return &adminRepo{db: db}
}

func (a *adminRepo) Create(ctx context.Context, req *repo.Admin) (*repo.Admin, error) {
	if err := a.db.WithContext(ctx).Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}


func (a *adminRepo) GetAll(ctx context.Context, status, firstname, lastname, email, phonenumber, username string, page, limit int) ([]repo.Admin, int, int, error) {
	var (
		admins     []repo.Admin
		total      int64
		totalPages int
	)

	dbQuery := a.db.WithContext(ctx).Model(&repo.Admin{})

	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if firstname != "" {
		dbQuery = dbQuery.Where("first_name LIKE ?", "%"+firstname+"%")
	}
	if lastname != "" {
		dbQuery = dbQuery.Where("last_name LIKE ?", "%"+lastname+"%")
	}
	if email != "" {
		dbQuery = dbQuery.Where("email LIKE ?", "%"+email+"%")
	}
	if phonenumber != "" {
		dbQuery = dbQuery.Where("phone_number LIKE ?", "%"+phonenumber+"%")
	}
	if username != "" {
		dbQuery = dbQuery.Where("username LIKE ?", "%"+username+"%")
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}
	totalPages = int(math.Ceil(float64(total) / float64(limit)))

	if err := dbQuery.Order("id ASC").Limit(limit).Offset((page - 1) * limit).Find(&admins).Error; err != nil {
		return nil, 0, 0, err
	}

	return admins, page, totalPages, nil
}

func (a *adminRepo) GetByUsername(ctx context.Context, username string) (*repo.Admin, error) {
	var admin repo.Admin

	if err := a.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepo) GetById(ctx context.Context, id int) (*repo.Admin, error) {
	var admin repo.Admin
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepo) UpdatePassword(ctx context.Context, username, password string) error {
	result := a.db.WithContext(ctx).Model(&repo.Admin{}).
		Where("username = ?", username).
		Update("password_hash", password)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *adminRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	result := a.db.WithContext(ctx).Model(&repo.Admin{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *adminRepo) Update(ctx context.Context, id int, req *repo.UpdateAdmin) error {
	result := a.db.WithContext(ctx).Model(&repo.Admin{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"first_name":    req.FirstName,
			"last_name":     req.LastName,
			"email":         req.Email,
			"phone_number":  req.PhoneNumber,
			"username":      req.Username,
			"password_hash": req.PasswordHash,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *adminRepo) DeleteById(ctx context.Context, id int) error {
	result := a.db.WithContext(ctx).Where("id = ?", id).Delete(&repo.Admin{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *adminRepo) DeleteByUsername(ctx context.Context, username string) error {
	result := a.db.WithContext(ctx).Where("username = ?", username).Delete(&repo.Admin{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

