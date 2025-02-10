package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type tokenRepo struct {
	db *sql.DB
}

func NewTokenStorage(db *sql.DB) repo.TokenStorageI {
	return &tokenRepo{
		db: db,
	}
}

func (t *tokenRepo) Create(ctx context.Context, req *repo.CreateToken) (*repo.CreateToken, error) {
	query := `
		INSERT INTO tokens (
			admin_id,
			token
		) VALUES (?, ?)
	`
	_, err := t.db.Exec(query, req.AdminId, req.Token)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (t *tokenRepo) GetByAdminId(ctx context.Context, id int) ([]repo.Token, error) {
	query := `
		SELECT 
			id, admin_id, 
			token, auth_time
		FROM tokens 
		WHERE admin_id = ?
	`

	rows, err := t.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []repo.Token

	for rows.Next() {
		var token repo.Token
		var auth_time string

		err := rows.Scan(
			&token.Id,
			&token.AdminId,
			&token.Token,
			&auth_time,
		)
		if err != nil {
			return nil, err
		}

		token.AuthTime, err = time.Parse("2006-01-02 15:04:05", auth_time)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	// Agar hech qanday token topilmasa, bo'sh slice qaytaramiz (`nil` emas)
	if len(tokens) == 0 {
		return []repo.Token{}, nil
	}

	return tokens, nil
}

func (t *tokenRepo) Delete(ctx context.Context, id int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM tokens WHERE id = ?`
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
