package storage

import (
	"github.com/sattorovshoxrux3009/Restourants_back/storage/mysql"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/gorm"
)

type StorageI interface {
	SuperAdmin() repo.SuperAdminStorageI
	Admin() repo.AdminStorageI
	Token() repo.TokenStorageI
	Restaurants() repo.RestaurantsI
	AdminRestaurantsLimit() repo.AdminRestaurantsLimitI
	Menu() repo.MenuI
	EventPrices() repo.EventPricesI
}
type StoragePg struct {
	superAdminRepo            repo.SuperAdminStorageI
	adminRepo                 repo.AdminStorageI
	tokenRepo                 repo.TokenStorageI
	restaurantsRepo           repo.RestaurantsI
	adminRestaurantsLimitRepo repo.AdminRestaurantsLimitI
	menuRepo                  repo.MenuI
	eventPricesRepo           repo.EventPricesI
}

func NewStorage(mysqlConn *gorm.DB) StorageI {
	return &StoragePg{
		superAdminRepo:            mysql.NewSuperAdminStorage(mysqlConn),
		adminRepo:                 mysql.NewAdminStorage(mysqlConn),
		tokenRepo:                 mysql.NewTokenStorage(mysqlConn),
		restaurantsRepo:           mysql.NewRestaurantsStorage(mysqlConn),
		adminRestaurantsLimitRepo: mysql.NewAdminRestaurantsLimitStorage(mysqlConn),
		menuRepo:                  mysql.NewMenuStorage(mysqlConn),
		eventPricesRepo:           mysql.NewEventPricesStorage(mysqlConn),
	}
}
func (s *StoragePg) SuperAdmin() repo.SuperAdminStorageI {
	return s.superAdminRepo
}
func (s *StoragePg) Admin() repo.AdminStorageI {
	return s.adminRepo
}
func (s *StoragePg) Token() repo.TokenStorageI {
	return s.tokenRepo
}
func (s *StoragePg) Restaurants() repo.RestaurantsI {
	return s.restaurantsRepo
}
func (s *StoragePg) AdminRestaurantsLimit() repo.AdminRestaurantsLimitI {
	return s.adminRestaurantsLimitRepo
}
func (s *StoragePg) Menu() repo.MenuI {
	return s.menuRepo
}
func (s *StoragePg) EventPrices() repo.EventPricesI {
	return s.eventPricesRepo
}
