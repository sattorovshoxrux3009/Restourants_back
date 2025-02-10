package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type restourantsRepo struct {
	db *sql.DB
}

func NewRestourantsStorage(db *sql.DB) repo.RestourantsI {
	return &restourantsRepo{
		db: db,
	}
}

func (r *restourantsRepo) Create(ctx context.Context, req *repo.CreateRestourant) (*repo.CreateRestourant, error) {
	query := `
		INSERT INTO restourants (
			name, address,
			latitude, longitude,
			phone_number, email,
			capacity, owner_id,
			opening_hours, image_url,
			description, alcohol_permission
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
		query, req.Name,
		req.Address, req.Latitude,
		req.Longitude, req.PhoneNumber,
		req.Email, req.Capacity,
		req.OwnerID, req.OpeningHours,
		req.ImageURL, req.Description,
		req.AlcoholPermission,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *restourantsRepo) GetByOwnerId(ctx context.Context, id int) ([]repo.Restaurant, error) {
	query := `
		SELECT 
			id, name, address,
			latitude, longitude,
			phone_number, email,
			capacity, owner_id,
			opening_hours, image_url,
			description, alcohol_permission,
			created_at, updated_at
		FROM restourants 
		WHERE owner_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restourants []repo.Restaurant

	for rows.Next() {
		var restourant repo.Restaurant
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&restourant.Id,
			&restourant.Name,
			&restourant.Address,
			&restourant.Latitude,
			&restourant.Longitude,
			&restourant.PhoneNumber,
			&restourant.Email,
			&restourant.Capacity,
			&restourant.OwnerID,
			&restourant.OpeningHours,
			&restourant.ImageURL,
			&restourant.Description,
			&restourant.AlcoholPermission,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Time parsing
		restourant.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		restourant.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, err
		}

		restourants = append(restourants, restourant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return restourants, nil
}

func (r *restourantsRepo) GetById(ctx context.Context, id int) (*repo.Restaurant, error) {
	query := `
		SELECT 
			id, name, address,
			latitude, longitude,
			phone_number, email,
			capacity, owner_id,
			opening_hours, image_url,
			description, alcohol_permission,
			created_at, updated_at
		FROM restourants 
		WHERE id = ?
	`
	var restourant repo.Restaurant
	var createdAtStr, updatedAtStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&restourant.Id,
		&restourant.Name,
		&restourant.Address,
		&restourant.Latitude,
		&restourant.Longitude,
		&restourant.PhoneNumber,
		&restourant.Email,
		&restourant.Capacity,
		&restourant.OwnerID,
		&restourant.OpeningHours,
		&restourant.ImageURL,
		&restourant.Description,
		&restourant.AlcoholPermission,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	restourant.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	restourant.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, err
	}
	return &restourant, nil
}

func (r *restourantsRepo) Update(ctx context.Context, id int, req *repo.UpdateRestourant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `
		UPDATE restourants SET 
			name = ?,
			address = ?,
			latitude = ?,
			longitude = ?,
			phone_number = ?,
			email = ?,
			capacity = ?,
			owner_id = ?,
			opening_hours = ?,
			image_url = ?,
			description = ?,
			alcohol_permission = ? 
		WHERE id = ?
	`
	_, err = tx.ExecContext(
		ctx, query, req.Name,
		req.Address, req.Latitude,
		req.Longitude, req.PhoneNumber,
		req.Email, req.Capacity, req.OwnerID,
		req.OpeningHours, req.ImageURL,
		req.Description, req.AlcoholPermission, id,
	)
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

func (r *restourantsRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM restourants WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, id)
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
