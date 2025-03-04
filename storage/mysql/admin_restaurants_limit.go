package mysql

import (
	"context"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type adminRestaurantsLimitRepo struct {
	db *gorm.DB
}

func NewAdminRestaurantsLimitStorage(db *gorm.DB) repo.AdminRestaurantsLimitI {
	return &adminRestaurantsLimitRepo{
		db: db,
	}
}

// **Create**
func (al *adminRestaurantsLimitRepo) Create(ctx context.Context, req *repo.AdminRestaurantLimit) (*repo.AdminRestaurantLimit, error) {
	err := al.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

// **GetByAdminId**
func (al *adminRestaurantsLimitRepo) GetByAdminId(ctx context.Context, id int) (*repo.AdminRestaurantLimit, error) {
	var limit repo.AdminRestaurantLimit
	err := al.db.WithContext(ctx).Where("admin_id = ?", id).First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &limit, nil
}

// **Update**
func (al *adminRestaurantsLimitRepo) Update(ctx context.Context, req *repo.AdminRestaurantLimit) error {
	err := al.db.WithContext(ctx).
		Model(&repo.AdminRestaurantLimit{}).
		Where("admin_id = ?", req.AdminId).
		Updates(map[string]interface{}{
			"max_restaurants": req.MaxRestaurants,
		}).Error
	return err
}
