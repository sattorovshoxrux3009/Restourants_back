package mysql

import (
	"context"
	"database/sql"

	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
)

type eventPricesRepo struct {
	db *sql.DB
}

func NewEventPricesStorage(db *sql.DB) repo.EventPricesI {
	return &eventPricesRepo{
		db: db,
	}
}

func (e *eventPricesRepo) Create(ctx context.Context, req *repo.CreateEventPrices) (*repo.CreateEventPrices, error) {
	query := `
		INSERT INTO event_prices (
			restaurant_id,
			event_type, table_price,
			waiter_price, max_guests,
			table_seats, max_waiters,
			alcohol_permission
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := e.db.Exec(
		query, req.RestaurantId,
		req.EventType, req.TablePrice,
		req.WaiterPrice, req.MaxGuests,
		req.TableSeats, req.MaxWaiters,
		req.AlcoholPermission,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (e *eventPricesRepo) GetAll(ctx context.Context, restaurantID, eventType string, page, limit int) ([]*repo.EventPrices, int, int, error) {
	query := `SELECT id, restaurant_id, event_type, table_price, waiter_price, 
		       max_guests, table_seats, max_waiters, alcohol_permission 
		FROM event_prices WHERE 1=1`
	args := []interface{}{}

	// Agar filtrlar mavjud bo‘lsa, ularni qo‘shamiz
	if restaurantID != "" {
		query += " AND restaurant_id = ?"
		args = append(args, restaurantID)
	}
	if eventType != "" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}

	// Tartiblash (restaurant_id bo‘yicha o‘sish tartibida)
	query += " ORDER BY restaurant_id ASC"

	// Pagination
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := e.db.Query(query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var eventPrices []*repo.EventPrices
	for rows.Next() {
		var event repo.EventPrices
		if err := rows.Scan(
			&event.Id, &event.RestaurantId, &event.EventType, &event.TablePrice,
			&event.WaiterPrice, &event.MaxGuests, &event.TableSeats,
			&event.MaxWaiters, &event.AlcoholPermission,
		); err != nil {
			return nil, 0, 0, err
		}
		eventPrices = append(eventPrices, &event)
	}

	// Jami sahifalar sonini hisoblash
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM event_prices WHERE 1=1"
	countArgs := []interface{}{}

	if restaurantID != "" {
		countQuery += " AND restaurant_id = ?"
		countArgs = append(countArgs, restaurantID)
	}
	if eventType != "" {
		countQuery += " AND event_type = ?"
		countArgs = append(countArgs, eventType)
	}

	err = e.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := (totalCount + limit - 1) / limit // Jami sahifalar

	if len(eventPrices) == 0 {
		return []*repo.EventPrices{}, page, totalPages, nil
	}

	return eventPrices, page, totalPages, nil
}

func (e *eventPricesRepo) GetByRestaurantID(ctx context.Context, restaurantID int) ([]*repo.EventPrices, error) {
	query := `
		SELECT id, restaurant_id, event_type, table_price, waiter_price, 
		       max_guests, table_seats, max_waiters, alcohol_permission 
		FROM event_prices WHERE restaurant_id = ?
	`

	rows, err := e.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventPrices []*repo.EventPrices
	for rows.Next() {
		var event repo.EventPrices
		if err := rows.Scan(
			&event.Id, &event.RestaurantId, &event.EventType, &event.TablePrice,
			&event.WaiterPrice, &event.MaxGuests, &event.TableSeats,
			&event.MaxWaiters, &event.AlcoholPermission,
		); err != nil {
			return nil, err
		}
		eventPrices = append(eventPrices, &event)
	}

	if len(eventPrices) == 0 {
		return nil, nil
	}

	return eventPrices, nil
}

func (e *eventPricesRepo) GetByID(ctx context.Context, id int) (*repo.EventPrices, error) {
	query := `
		SELECT id, restaurant_id, event_type, table_price, waiter_price, 
		       max_guests, table_seats, max_waiters, alcohol_permission 
		FROM event_prices WHERE id = ?
	`

	var event repo.EventPrices
	err := e.db.QueryRow(query, id).Scan(
		&event.Id, &event.RestaurantId, &event.EventType, &event.TablePrice,
		&event.WaiterPrice, &event.MaxGuests, &event.TableSeats,
		&event.MaxWaiters, &event.AlcoholPermission,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *eventPricesRepo) Update(ctx context.Context, event *repo.EventPrices) error {
	tx, err := e.db.Begin() // Transaction boshlash
	if err != nil {
		return err
	}

	query := `
		UPDATE event_prices SET 
			table_price = ?, 
			waiter_price = ?, max_guests = ?, table_seats = ?, 
			max_waiters = ?, alcohol_permission = ? 
		WHERE id = ?
	`

	_, err = tx.Exec(query,
		event.TablePrice, event.WaiterPrice, event.MaxGuests, event.TableSeats,
		event.MaxWaiters, event.AlcoholPermission, event.Id,
	)

	if err != nil {
		tx.Rollback() // Xatolik bo‘lsa transactionni bekor qilish
		return err
	}

	err = tx.Commit() // O‘zgarishlarni tasdiqlash
	if err != nil {
		return err
	}

	return nil
}

func (e *eventPricesRepo) Delete(ctx context.Context, id int) error {
	tx, err := e.db.Begin() // Transaction boshlash
	if err != nil {
		return err
	}

	query := `DELETE FROM event_prices WHERE id = ?`

	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback() // Xatolik bo‘lsa transactionni bekor qilish
		return err
	}

	err = tx.Commit() // O‘zgarishlarni tasdiqlash
	if err != nil {
		return err
	}

	return nil
}
