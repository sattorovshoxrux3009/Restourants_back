package mysql

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type restaurantsRepo struct {
	db *sql.DB
}

func NewRestaurantsStorage(db *sql.DB) repo.RestaurantsI {
	return &restaurantsRepo{
		db: db,
	}
}

func (r *restaurantsRepo) Create(ctx context.Context, req *repo.CreateRestaurant) (*repo.CreateRestaurant, error) {
	query := `
		INSERT INTO restaurants (
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

// for users, unlocked
func (r *restaurantsRepo) GetAll(ctx context.Context, name, address, capacity, adlcohol_permission string, page, limit int) ([]repo.Restaurant, int, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM restaurants WHERE status = 'active'`
	var args []interface{}

	// Qidiruv parametrlarini qo‘shish
	if name != "" {
		countQuery += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if address != "" {
		countQuery += " AND address LIKE ?"
		args = append(args, "%"+address+"%")
	}

	if capacity != "" {
		countQuery += " AND capacity LIKE ?"
		args = append(args, "%"+capacity+"%")
	}

	if adlcohol_permission != "" {
		countQuery += " AND adlcohol_permission LIKE ?"
		args = append(args, "%"+adlcohol_permission+"%")
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Adminlarni olish uchun so‘rov
	query := `SELECT id, name, address, latitude, longitude, phone_number, email, capacity, owner_id, opening_hours, image_url, description, alcohol_permission 
              FROM restaurants WHERE status = 'active'`
	args = nil // Fresh args list

	// Qidiruv parametrlarini qo‘shish
	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if address != "" {
		query += " AND address LIKE ?"
		args = append(args, "%"+address+"%")
	}

	if capacity != "" {
		query += " AND capacity LIKE ?"
		args = append(args, "%"+capacity+"%")
	}

	if adlcohol_permission != "" {
		query += " AND adlcohol_permission LIKE ?"
		args = append(args, "%"+adlcohol_permission+"%")
	}

	// Sahifani tartiblaymiz va limit qo‘shamiz
	query += " ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	// So‘rovni bajarish
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var restaurants []repo.Restaurant

	// Natijalarni yig‘ish
	for rows.Next() {
		var restaurant repo.Restaurant

		err := rows.Scan(
			&restaurant.Id,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&restaurant.PhoneNumber,
			&restaurant.Email,
			&restaurant.Capacity,
			&restaurant.OwnerID,
			&restaurant.OpeningHours,
			&restaurant.ImageURL,
			&restaurant.Description,
			&restaurant.AlcoholPermission,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return restaurants, page, totalPages, nil
}

// for super admin
func (r *restaurantsRepo) GetSall(ctx context.Context, status, name, address, capacity, adlcohol_permission string, page, limit int) ([]repo.Restaurant, int, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM restaurants`
	var args []interface{}

	// Qidiruv parametrlarini qo‘shish
	if status != "" {
		countQuery += " WHERE status = ?"
		args = append(args, status)
	}

	if name != "" {
		if len(args) > 0 {
			countQuery += " AND name LIKE ?"
		} else {
			countQuery += " WHERE name LIKE ?"
		}
		args = append(args, "%"+name+"%")
	}

	if address != "" {
		if len(args) > 0 {
			countQuery += " AND address LIKE ?"
		} else {
			countQuery += " WHERE address LIKE ?"
		}
		args = append(args, "%"+address+"%")
	}

	if capacity != "" {
		if len(args) > 0 {
			countQuery += " AND capacity LIKE ?"
		} else {
			countQuery += " WHERE capacity LIKE ?"
		}
		args = append(args, "%"+capacity+"%")
	}

	if adlcohol_permission != "" {
		if len(args) > 0 {
			countQuery += " AND adlcohol_permission LIKE ?"
		} else {
			countQuery += " WHERE adlcohol_permission LIKE ?"
		}
		args = append(args, "%"+adlcohol_permission+"%")
	}

	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Adminlarni olish uchun so‘rov
	query := `SELECT id, name, address, latitude, longitude, phone_number, email, capacity, owner_id, opening_hours, image_url, description, alcohol_permission 
	          FROM restaurants`
	args = nil // Fresh args list

	// Qidiruv parametrlarini qo‘shish
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	if name != "" {
		if len(args) > 0 {
			query += " AND name LIKE ?"
		} else {
			query += " WHERE name LIKE ?"
		}
		args = append(args, "%"+name+"%")
	}

	if address != "" {
		if len(args) > 0 {
			query += " AND address LIKE ?"
		} else {
			query += " WHERE address LIKE ?"
		}
		args = append(args, "%"+address+"%")
	}

	if capacity != "" {
		if len(args) > 0 {
			query += " AND capacity LIKE ?"
		} else {
			query += " WHERE capacity LIKE ?"
		}
		args = append(args, "%"+capacity+"%")
	}

	if adlcohol_permission != "" {
		if len(args) > 0 {
			query += " AND adlcohol_permission LIKE ?"
		} else {
			query += " WHERE adlcohol_permission LIKE ?"
		}
		args = append(args, "%"+adlcohol_permission+"%")
	}

	// Sahifani tartiblaymiz va limit qo‘shamiz
	query += " ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	// So‘rovni bajarish
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var restaurants []repo.Restaurant

	// Natijalarni yig‘ish
	for rows.Next() {
		var restaurant repo.Restaurant

		err := rows.Scan(
			&restaurant.Id,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&restaurant.PhoneNumber,
			&restaurant.Email,
			&restaurant.Capacity,
			&restaurant.OwnerID,
			&restaurant.OpeningHours,
			&restaurant.ImageURL,
			&restaurant.Description,
			&restaurant.AlcoholPermission,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return restaurants, page, totalPages, nil
}

func (r *restaurantsRepo) GetByOwnerId(ctx context.Context, id int) ([]repo.Restaurant, error) {
	query := `
		SELECT 
			id, name, address,
			latitude, longitude,
			phone_number, email,
			capacity, owner_id,
			opening_hours, image_url,
			description, alcohol_permission,
			created_at, updated_at,
			status
		FROM restaurants 
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
			&restourant.Status,
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

func (r *restaurantsRepo) GetById(ctx context.Context, id int) (*repo.Restaurant, error) {
	query := `
		SELECT 
			id, name, address,
			latitude, longitude,
			phone_number, email,
			capacity, owner_id,
			opening_hours, image_url,
			description, alcohol_permission,
			created_at, updated_at,
			status
		FROM restaurants 
		WHERE id = ?
	`
	var restaurant repo.Restaurant
	var createdAtStr, updatedAtStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&restaurant.Id,
		&restaurant.Name,
		&restaurant.Address,
		&restaurant.Latitude,
		&restaurant.Longitude,
		&restaurant.PhoneNumber,
		&restaurant.Email,
		&restaurant.Capacity,
		&restaurant.OwnerID,
		&restaurant.OpeningHours,
		&restaurant.ImageURL,
		&restaurant.Description,
		&restaurant.AlcoholPermission,
		&createdAtStr,
		&updatedAtStr,
		&restaurant.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	restaurant.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	restaurant.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantsRepo) Update(ctx context.Context, id int, req *repo.UpdateRestaurant) error {
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
		UPDATE restaurants SET 
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

func (r *restaurantsRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM restaurants WHERE id = ?`
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
