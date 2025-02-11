package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type adminRestaurantsLimitRepo struct {
	db *sql.DB
}

func NewAdminRestaurantsLimitStorage(db *sql.DB) repo.AdminRestaurantsLimitI {
	return &adminRestaurantsLimitRepo{
		db: db,
	}
}

func (al *adminRestaurantsLimitRepo) Create(ctx context.Context, req *repo.CreateAdminRestaurantsLimit) (*repo.CreateAdminRestaurantsLimit, error) {
	query := `
		INSERT INTO admin_restaurant_limits (
			admin_id,
			max_restaurants
		) VALUES (?, ?)
	`
	_, err := al.db.Exec(query, req.AdminId, req.MaxRestaurants)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (al *adminRestaurantsLimitRepo) GetByAdminId(ctx context.Context, id int) (*repo.AdminRestaurantsLimit, error) {
	query := `
		SELECT 
			id,
			admin_id,
			max_restaurants,
			created_at, updated_at
		FROM admin_restaurant_limits 
		WHERE admin_id = ?
	`
	var limit repo.AdminRestaurantsLimit
	var createdAtStr, updatedAtStr string
	err := al.db.QueryRowContext(ctx, query, id).Scan(
		&limit.Id,
		&limit.AdminId,
		&limit.MaxRestaurants,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	limit.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	limit.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (al *adminRestaurantsLimitRepo) Update(ctx context.Context, req *repo.CreateAdminRestaurantsLimit) error {
	tx, err := al.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `
		UPDATE admin_restaurant_limits SET 
			admin_id = ?,
			max_restaurants = ?
		WHERE admin_id = ?
	`
	_, err = tx.ExecContext(ctx, query, req.AdminId, req.MaxRestaurants, req.AdminId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
