package mysql

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminStorage(db *sql.DB) repo.AdminStorageI {
	return &adminRepo{
		db: db,
	}
}

func (a *adminRepo) Create(ctx context.Context, req *repo.CreateAdmin) (*repo.CreateAdmin, error) {
	query := `
		INSERT INTO admins (
			first_name,
			last_name,
			email,
			phone_number,
			username,
			password_hash
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := a.db.Exec(query, req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Username, req.PasswordHash)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (a *adminRepo) GetAll(ctx context.Context, status, firstname, lastname, email, phonenumber, username string, page, limit int) ([]repo.Admin, int, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM admins`
	var args []interface{}

	// Qidiruv parametrlarini qo‘shish
	if status != "" {
		countQuery += " WHERE status = ?"
		args = append(args, status)
	}

	if firstname != "" {
		if len(args) > 0 {
			countQuery += " AND first_name LIKE ?"
		} else {
			countQuery += " WHERE first_name LIKE ?"
		}
		args = append(args, "%"+firstname+"%")
	}

	if lastname != "" {
		if len(args) > 0 {
			countQuery += " AND last_name LIKE ?"
		} else {
			countQuery += " WHERE last_name LIKE ?"
		}
		args = append(args, "%"+lastname+"%")
	}

	if email != "" {
		if len(args) > 0 {
			countQuery += " AND email LIKE ?"
		} else {
			countQuery += " WHERE email LIKE ?"
		}
		args = append(args, "%"+email+"%")
	}

	if phonenumber != "" {
		if len(args) > 0 {
			countQuery += " AND phone_number LIKE ?"
		} else {
			countQuery += " WHERE phone_number LIKE ?"
		}
		args = append(args, "%"+phonenumber+"%")
	}

	if username != "" {
		if len(args) > 0 {
			countQuery += " AND username LIKE ?"
		} else {
			countQuery += " WHERE username LIKE ?"
		}
		args = append(args, "%"+username+"%")
	}

	// Umumiy sonni olish
	err := a.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Adminlarni olish uchun so‘rov
	query := `SELECT id, first_name, last_name, email, phone_number, username, password_hash, created_at, status 
	          FROM admins`
	args = nil // Fresh args list

	// Qidiruv parametrlarini qo‘shish
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	if firstname != "" {
		if len(args) > 0 {
			query += " AND first_name LIKE ?"
		} else {
			query += " WHERE first_name LIKE ?"
		}
		args = append(args, "%"+firstname+"%")
	}

	if lastname != "" {
		if len(args) > 0 {
			query += " AND last_name LIKE ?"
		} else {
			query += " WHERE last_name LIKE ?"
		}
		args = append(args, "%"+lastname+"%")
	}

	if email != "" {
		if len(args) > 0 {
			query += " AND email LIKE ?"
		} else {
			query += " WHERE email LIKE ?"
		}
		args = append(args, "%"+email+"%")
	}

	if phonenumber != "" {
		if len(args) > 0 {
			query += " AND phone_number LIKE ?"
		} else {
			query += " WHERE phone_number LIKE ?"
		}
		args = append(args, "%"+phonenumber+"%")
	}

	if username != "" {
		if len(args) > 0 {
			query += " AND username LIKE ?"
		} else {
			query += " WHERE username LIKE ?"
		}
		args = append(args, "%"+username+"%")
	}

	// Sahifani tartiblaymiz va limit qo‘shamiz
	query += " ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	// So‘rovni bajarish
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var admins []repo.Admin

	// Natijalarni yig‘ish
	for rows.Next() {
		var admin repo.Admin
		var createdAtStr string

		err := rows.Scan(
			&admin.Id,
			&admin.FirstName,
			&admin.LastName,
			&admin.Email,
			&admin.PhoneNumber,
			&admin.Username,
			&admin.PasswordHash,
			&createdAtStr,
			&admin.Status,
		)
		if err != nil {
			return nil, 0, 0, err
		}

		admin.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, 0, 0, err
		}

		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return admins, page, totalPages, nil
}

func (a *adminRepo) GetByUsername(ctx context.Context, username string) (*repo.Admin, error) {
	query := `
		SELECT 
			id,	first_name, 
			last_name, 
			email,phone_number, 
			username, 
			password_hash, 
			created_at, status 
		FROM admins 
		WHERE username = ?
	`
	var admin repo.Admin
	var createdAtStr string
	err := a.db.QueryRowContext(ctx, query, username).Scan(
		&admin.Id,
		&admin.FirstName,
		&admin.LastName,
		&admin.Email,
		&admin.PhoneNumber,
		&admin.Username,
		&admin.PasswordHash,
		&createdAtStr,
		&admin.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	admin.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepo) GetById(ctx context.Context, id int) (*repo.Admin, error) {
	query := `
		SELECT 
			id,	first_name, 
			last_name, 
			email,phone_number, 
			username, 
			password_hash, 
			created_at, status 
		FROM admins 
		WHERE id = ?
	`
	var admin repo.Admin
	var createdAtStr string
	err := a.db.QueryRowContext(ctx, query, id).Scan(
		&admin.Id,
		&admin.FirstName,
		&admin.LastName,
		&admin.Email,
		&admin.PhoneNumber,
		&admin.Username,
		&admin.PasswordHash,
		&createdAtStr,
		&admin.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	admin.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepo) UpdatePassword(ctx context.Context, username, password string) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE admins SET password_hash = ? WHERE username = ?`
	_, err = tx.ExecContext(ctx, query, password, username)
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

func (a *adminRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE admins SET status = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, status, id)
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

func (a *adminRepo) Update(ctx context.Context, id int, req *repo.UpdateAdmin) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `
		UPDATE admins SET 
			first_name = ?,
			last_name = ?,
			email = ?,
			phone_number = ?,
			username = ?,
			password_hash = ? 
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, query, req.FirstName, req.LastName, req.Email, req.PhoneNumber, req.Username, req.PasswordHash, id)
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

func (a *adminRepo) DeleteById(ctx context.Context, id int) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM admins WHERE id = ?`
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

func (a *adminRepo) DeleteByUsername(ctx context.Context, username string) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM admins WHERE username = ?`
	_, err = tx.ExecContext(ctx, query, username)
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
