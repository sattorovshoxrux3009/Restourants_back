package mysql

import (
	"context"
	"errors"
	"math"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type restaurantsRepo struct {
	db *gorm.DB
}

func NewRestaurantsStorage(db *gorm.DB) repo.RestaurantsI {
	return &restaurantsRepo{
		db: db,
	}
}

func (r *restaurantsRepo) Create(ctx context.Context, req *repo.Restaurant) (*repo.Restaurant, error) {
	restaurant := repo.Restaurant{
		Name:              req.Name,
		Address:           req.Address,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		PhoneNumber:       req.PhoneNumber,
		Email:             req.Email,
		Capacity:          req.Capacity,
		OwnerId:           uint(req.OwnerId),
		OpeningHours:      req.OpeningHours,
		ImageURL:          req.ImageURL,
		Description:       req.Description,
		AlcoholPermission: req.AlcoholPermission,
	}

	if err := r.db.Create(&restaurant).Error; err != nil {
		return nil, err
	}
	return req, nil
}

// for users, unlocked
func (r *restaurantsRepo) GetAll(ctx context.Context, name, address, capacity, alcoholPermission string, page, limit int) ([]repo.Restaurant, int, int, error) {
	var total int64
	var restaurants []repo.Restaurant
	var err error

	query := r.db.Model(&repo.Restaurant{}).Where("status = ?", "active")

	// Qidiruv parametrlarini qo‘shish
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if address != "" {
		query = query.Where("address LIKE ?", "%"+address+"%")
	}
	if capacity != "" {
		query = query.Where("capacity LIKE ?", "%"+capacity+"%")
	}
	if alcoholPermission != "" {
		query = query.Where("alcohol_permission LIKE ?", "%"+alcoholPermission+"%")
	}

	// So'rovni bajarish
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Qidiruv natijalarini olish
	err = query.Offset((page - 1) * limit).Limit(limit).Find(&restaurants).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return restaurants, page, totalPages, nil
}

// for super admin
func (r *restaurantsRepo) GetSall(ctx context.Context, status, phonenumber, email, ownerid, name, address, capacity, alcoholPermission string, page, limit int) ([]repo.Restaurant, int, int, error) {
	var total int64
	var restaurants []repo.Restaurant
	var err error

	query := r.db.Model(&repo.Restaurant{})

	// Qidiruv parametrlarini qo‘shish
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if phonenumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+phonenumber+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if ownerid != "" {
		query = query.Where("owner_id LIKE ?", "%"+ownerid+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if address != "" {
		query = query.Where("address LIKE ?", "%"+address+"%")
	}
	if capacity != "" {
		query = query.Where("capacity LIKE ?", "%"+capacity+"%")
	}
	if alcoholPermission != "" {
		query = query.Where("alcohol_permission LIKE ?", "%"+alcoholPermission+"%")
	}

	// So'rovni bajarish
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Qidiruv natijalarini olish
	err = query.Offset((page - 1) * limit).Limit(limit).Find(&restaurants).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return restaurants, page, totalPages, nil
}

func (r *restaurantsRepo) GetByOwnerId(ctx context.Context, ownerID int, name string, limit int) ([]repo.Restaurant, error) {
	var restaurants []repo.Restaurant

	query := r.db.WithContext(ctx).Where("owner_id = ?", ownerID)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	result := query.Limit(limit).Find(&restaurants)
	return restaurants, result.Error
}

func (r *restaurantsRepo) GetById(ctx context.Context, id int) (*repo.Restaurant, error) {
	var restaurant repo.Restaurant
	result := r.db.WithContext(ctx).First(&restaurant, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("restaurant not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *restaurantsRepo) Update(ctx context.Context, id int, req *repo.Restaurant) error {

	result := r.db.WithContext(ctx).Model(&repo.Restaurant{}).Where("id = ?", id).Updates(req)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("restaurant not found")
	}

	return nil
}

func (r *restaurantsRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	result := r.db.WithContext(ctx).Model(&repo.Restaurant{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("restaurant not found")
	}

	return nil
}

func (r *restaurantsRepo) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&repo.Restaurant{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("restaurant not found")
	}

	return nil
}
