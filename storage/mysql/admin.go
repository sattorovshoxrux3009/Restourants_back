package mysql

import (
	"context"
	"database/sql"
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

func (a *adminRepo) GetAll(ctx context.Context) ([]repo.Admin, error) {
	query := `SELECT id, first_name, last_name, email, phone_number, username, password_hash, created_at, status FROM admins`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []repo.Admin

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
			return nil, err
		}
		admin.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
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
