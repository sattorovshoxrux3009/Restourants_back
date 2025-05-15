package mysql

import (
	"context"
	"math"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type menuRepo struct {
	db *gorm.DB
}

func NewMenuStorage(db *gorm.DB) repo.MenuI {
	return &menuRepo{db: db}
}

func (m *menuRepo) Create(ctx context.Context, req *repo.Menu) (*repo.Menu, error) {
	err := m.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (m *menuRepo) GetAll(ctx context.Context, name, category string, page, limit int) ([]repo.Menu, int, int, error) {
	var menus []repo.Menu
	var total int64

	query := m.db.WithContext(ctx).Model(&repo.Menu{}).
		Joins("JOIN restaurant r ON menu.restaurant_id = r.id").
		Where("r.status = ?", "active")

	if name != "" {
		query = query.Where("menu.name LIKE ?", "%"+name+"%")
	}

	if category != "" {
		query = query.Where("menu.category LIKE ?", "%"+category+"%")
	}

	query.Count(&total)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	err := query.Order("menu.id ASC").Limit(limit).Offset((page - 1) * limit).
		Find(&menus).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return menus, page, totalPages, nil
}

func (m *menuRepo) GetSAll(ctx context.Context, name, category string, restaurantID, page, limit int) ([]repo.Menu, int, int, error) {
	var menus []repo.Menu
	var total int64

	query := m.db.WithContext(ctx).Model(&repo.Menu{}).
		Select("menu.*, r.status AS restaurant_status"). // restaurant_status nomi bilan olamiz
		Joins("JOIN restaurant r ON menu.restaurant_id = r.id")

	if name != "" {
		query = query.Where("menu.name LIKE ?", "%"+name+"%")
	}

	if category != "" {
		query = query.Where("menu.category LIKE ?", "%"+category+"%")
	}

	if restaurantID != 0 {
		query = query.Where("menu.restaurant_id = ?", restaurantID)
	}

	// Jami yozuvlar sonini hisoblash
	query.Count(&total)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Query bajarish
	err := query.Order("menu.id ASC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&menus).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return menus, page, totalPages, nil
}

func (m *menuRepo) GetSAllByRestaurants(ctx context.Context, name, category string, restaurantIDs []uint, page, limit int) ([]repo.Menu, int, int, error) {
	var menus []repo.Menu
	var total int64

	// Agar adminning restorani bo'lmasa, bo'sh array qaytaramiz
	if len(restaurantIDs) == 0 {
		return []repo.Menu{}, page, 0, nil
	}

	// Query yaratish
	query := m.db.WithContext(ctx).Model(&repo.Menu{}).Where("restaurant_id IN ?", restaurantIDs)

	// Filtrlarni qo‘shish
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Jami natijalarni sanash
	query.Count(&total)

	// Pagination (sahifalash) qilish
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&menus).Error; err != nil {
		return nil, 0, 0, err
	}

	// Jami sahifalar sonini hisoblash
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return menus, page, totalPages, nil
}

func (m *menuRepo) GetById(ctx context.Context, id int) (*repo.Menu, error) {
	var menu repo.Menu
	err := m.db.WithContext(ctx).First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (m *menuRepo) GetByRestaurantId(ctx context.Context, id int) ([]*repo.Menu, error) {
	var menus []*repo.Menu

	if err := m.db.WithContext(ctx).Where("restaurant_id = ?", id).Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *menuRepo) Update(ctx context.Context, id int, req *repo.Menu) (*repo.Menu, error) {
	result := m.db.WithContext(ctx).Model(&repo.Menu{}).Where("id = ?", id).Updates(req)

	if result.Error != nil {
		return nil, result.Error
	}

	return req, nil
}

func (m *menuRepo) Delete(ctx context.Context, id int) error {
	if err := m.db.WithContext(ctx).Where("id = ?", id).Delete(&repo.Menu{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *menuRepo) SearchByName(ctx context.Context, nameQuery string, page int, limit int) ([]repo.MenuShort, int, error) {
	var results []repo.MenuShort
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	likePattern := nameQuery + "%"

	// 1. Jami yozuvlar sonini hisoblash
	var totalCount int64
	err := r.db.WithContext(ctx).
		Table("menu").
		Joins("JOIN restaurant ON menu.restaurant_id = restaurant.id").
		Where("menu.name LIKE ?", likePattern).
		Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. Menularni olish
	err = r.db.WithContext(ctx).
		Table("menu").
		Select(`menu.name, menu.image_url, menu.price, menu.restaurant_id, restaurant.name AS restaurant_name, ? AS type`, "menu").
		Joins("JOIN restaurant ON menu.restaurant_id = restaurant.id").
		Where("menu.name LIKE ?", likePattern).
		Limit(limit).
		Offset(offset).
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return results, totalPages, nil
}
