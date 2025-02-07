package storage

import (
	"database/sql"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/mysql"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type StorageI interface {
	SuperAdmin() repo.SuperAdminStorageI
	Admin() repo.AdminStorageI
}
type storagePg struct {
	superAdminRepo repo.SuperAdminStorageI
	adminRepo      repo.AdminStorageI
}

func NewStorage(mysqlConn *sql.DB) StorageI {
	return &storagePg{
		superAdminRepo: mysql.NewSuperAdminStorage(mysqlConn),
		adminRepo:      mysql.NewAdminStorage(mysqlConn),
	}
}
func (s *storagePg) SuperAdmin() repo.SuperAdminStorageI {
	return s.superAdminRepo
}
func (s *storagePg) Admin() repo.AdminStorageI {
	return s.adminRepo
}
