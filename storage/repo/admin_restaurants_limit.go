package repo

import (
	"context"
	"time"
)

type AdminRestaurantsLimitI interface {
	Create(ctx context.Context, req *AdminRestaurantLimit) (*AdminRestaurantLimit, error)
	GetByAdminId(ctx context.Context, id int) (*AdminRestaurantLimit, error)
	Update(ctx context.Context, req *AdminRestaurantLimit) error
}

//	type AdminRestaurantsLimit struct {
//		Id             int
//		AdminId        int
//		MaxRestaurants int
//		CreatedAt      time.Time
//		UpdatedAt      time.Time
//	}
// type CreateAdminRestaurantsLimit struct {
// 	AdminId        int
// 	MaxRestaurants int
// }

// AdminRestaurantLimits jadvali
type AdminRestaurantLimit struct {
	Id             uint      `gorm:"primaryKey"`
	AdminId        uint      `gorm:"not null"`
	MaxRestaurants int       `gorm:"not null;default:1"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	Admin          Admin     `gorm:"foreignKey:AdminId;constraint:OnDelete:CASCADE"`
}
