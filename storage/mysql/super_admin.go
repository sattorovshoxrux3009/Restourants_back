package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type superAdminRepo struct {
	db *sql.DB
}

func NewSuperAdminStorage(db *sql.DB) repo.SuperAdminStorageI {
	return &superAdminRepo{
		db: db,
	}
}

func (s *superAdminRepo) Create(ctx context.Context, req *repo.CreateSuperAdmin) (*repo.CreateSuperAdmin, error) {
	query := `
		INSERT INTO super_admins (
			username,
			password
		) VALUES (?, ?)
	`
	_, err := s.db.Exec(query, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	return req, nil
}

func (s *superAdminRepo) GetByUsername(ctx context.Context, username string) (*repo.SuperAdmin, error) {
	query := `
		SELECT 
			id, username,
			password, created_at, 
			token, last_login
		FROM super_admins WHERE username=?
	`
	var admin repo.SuperAdmin
	var createdAt, last_login []byte
	var token sql.NullString
	err := s.db.QueryRow(query, username).Scan(
		&admin.Id,
		&admin.Username,
		&admin.Password,
		&createdAt,
		&token,
		&last_login,
	)
	if err != nil {
		return nil, err
	}
	if len(createdAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}
		admin.CreatedAt = parsedTime
	}
	if len(last_login) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(last_login))
		if err != nil {
			return nil, err
		}
		admin.LastLogin = parsedTime
	}
	if token.Valid {
		admin.Token = token.String
	} else {
		admin.Token = ""
	}
	return &admin, nil
}

func (s *superAdminRepo) GetToken(ctx context.Context, username string) (string, error) {
	query := `
		SELECT  
			token
		FROM super_admins WHERE username=?
	`
	var token sql.NullString
	err := s.db.QueryRow(query, username).Scan(
		&token,
	)
	if err != nil {
		return "", err
	}
	if token.Valid {
		return token.String, err
	} else {
		return "", err
	}
}

func (s *superAdminRepo) UpdatePassword(ctx context.Context, username, password string) error {
	tsx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE super_admins SET 
			password=?
		WHERE username=?
	`
	res, err := tsx.Exec(query, password, username)
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}
	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}
	if data == 0 {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return sql.ErrNoRows
	}
	err = tsx.Commit()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}

	return nil
}

func (s *superAdminRepo) UpdateToken(ctx context.Context, username, token string) error {
	tsx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE super_admins SET 
			token=?
		WHERE username=?
	`
	res, err := tsx.Exec(query, token, username)
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}
	data, err := res.RowsAffected()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}
	if data == 0 {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return sql.ErrNoRows
	}
	err = tsx.Commit()
	if err != nil {
		errRoll := tsx.Rollback()
		if errRoll != nil {
			return errRoll
		}
		return err
	}

	return nil
}
