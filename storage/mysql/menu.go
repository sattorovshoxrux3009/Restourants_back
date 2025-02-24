package mysql

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type menuRepo struct {
	db *sql.DB
}

func NewMenuStorage(db *sql.DB) repo.MenuI {
	return &menuRepo{
		db: db,
	}
}

func (m *menuRepo) Create(ctx context.Context, req *repo.CreateMenu) (*repo.CreateMenu, error) {
	query := `
		INSERT INTO menu (
			restaurant_id,
			name, description,
			price, category,
			image_url
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := m.db.Exec(
		query, req.RestaurantId,
		req.Name, req.Description,
		req.Price, req.Category,
		req.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// for users, unlocked
func (m *menuRepo) GetAll(ctx context.Context, name, category string, page, limit int) ([]repo.Menu, int, int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) 
        FROM menu m 
        JOIN restaurants r ON m.restaurant_id = r.id 
        WHERE r.status = 'active'
	`
	var args []interface{}

	if name != "" {
		countQuery += " AND m.name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		countQuery += " AND m.category LIKE ?"
		args = append(args, "%"+category+"%")
	}

	// Umumiy natijalarni olish
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Ma'lumotlarni olish uchun so‘rov
	query := `SELECT m.id, m.restaurant_id, m.name, m.description, m.price, m.category,m.image_url
              FROM menu m
              JOIN restaurants r ON m.restaurant_id = r.id
              WHERE r.status = 'active'`
	args = nil // Yangi argumentlar ro‘yxati

	// Filtrlarni qo‘shish
	if name != "" {
		query += " AND m.name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		query += " AND m.category LIKE ?"
		args = append(args, "%"+category+"%")
	}

	// Sahifalash va tartiblash
	query += " ORDER BY m.id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	// So‘rovni bajarish
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var menus []repo.Menu

	for rows.Next() {
		var menu repo.Menu
		err := rows.Scan(
			&menu.Id,
			&menu.RestaurantId,
			&menu.Name,
			&menu.Description,
			&menu.Price,
			&menu.Category,
			&menu.ImageURL,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		menus = append(menus, menu)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return menus, page, totalPages, nil
}

// for super admin
func (m *menuRepo) GetSAll(ctx context.Context, name, category string, restaurant_id, page, limit int) ([]repo.MenuWithStatus, int, int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) 
        FROM menu m 
        JOIN restaurants r ON m.restaurant_id = r.id
        WHERE 1=1
	`
	var args []interface{}

	if name != "" {
		countQuery += " AND m.name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		countQuery += " AND m.category LIKE ?"
		args = append(args, "%"+category+"%")
	}

	if restaurant_id != 0 {
		countQuery += " AND m.restaurant_id = ?"
		args = append(args, restaurant_id)
	}

	// Umumiy natijalarni olish
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, 0, err
	}

	// Nechta sahifa borligini hisoblash
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Ma'lumotlarni olish uchun so‘rov
	query := `
		SELECT m.id, m.restaurant_id, m.name, m.description, m.price, m.category, m.image_url, m.created_at, m.updated_at, 
		       r.status 
		FROM menu m
		JOIN restaurants r ON m.restaurant_id = r.id
		WHERE 1=1
	`
	args = nil // Yangi argumentlar ro‘yxati

	// Filtrlarni qo‘shish
	if name != "" {
		query += " AND m.name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		query += " AND m.category LIKE ?"
		args = append(args, "%"+category+"%")
	}

	if restaurant_id != 0 {
		query += " AND m.restaurant_id = ?"
		args = append(args, restaurant_id)
	}

	// Sahifalash va tartiblash
	query += " ORDER BY m.id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, (page-1)*limit)

	// So‘rovni bajarish
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var menus []repo.MenuWithStatus

	for rows.Next() {
		var menu repo.MenuWithStatus
		var createdAtStr, updatedAtStr, restaurantStatus string
		err := rows.Scan(
			&menu.Id,
			&menu.RestaurantId,
			&menu.Name,
			&menu.Description,
			&menu.Price,
			&menu.Category,
			&menu.ImageURL,
			&createdAtStr,
			&updatedAtStr,
			&restaurantStatus,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		menu.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		menu.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		// Restaurant statusi bo'yicha menyu statusini belgilash
		if restaurantStatus == "active" {
			menu.Status = "active"
		} else {
			menu.Status = "inactive"
		}

		menus = append(menus, menu)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return menus, page, totalPages, nil
}

func (m *menuRepo) GetById(ctx context.Context, id int) (*repo.Menu, error) {
	query := `SELECT id, restaurant_id, name, description, price, category, image_url, created_at, updated_at 
	          FROM menu 
	          WHERE id = ?`

	var menu repo.Menu
	var createdAtStr, updatedAtStr string

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&menu.Id,
		&menu.RestaurantId,
		&menu.Name,
		&menu.Description,
		&menu.Price,
		&menu.Category,
		&menu.ImageURL,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Vaqt formatini pars qilish
	menu.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	menu.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

	return &menu, nil
}

func (m *menuRepo) Update(ctx context.Context, id int, req *repo.CreateMenu) (*repo.CreateMenu, error) {
	// Tranzaksiyani boshlash
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}

	// Agar xatolik bo‘lsa, tranzaksiyani bekor qilish (rollback)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE menu 
	          SET name = ?, description = ?, price = ?, category = ?, image_url = ? 
	          WHERE id = ?`

	_, err = tx.Exec(query, req.Name, req.Description, req.Price, req.Category, req.ImageURL, id)
	if err != nil {
		return nil, err
	}

	// Agar hammasi yaxshi bo‘lsa, tranzaksiyani tasdiqlash (commit)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (m *menuRepo) Delete(ctx context.Context, id int) error {
	
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `DELETE FROM menu WHERE id = ?`

	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
