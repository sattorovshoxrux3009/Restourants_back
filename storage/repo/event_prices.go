package repo

import "context"

type EventPricesI interface {
	Create(ctx context.Context, req *CreateEventPrices) (*CreateEventPrices, error)
	GetAll(ctx context.Context, restaurantID, eventType string, page, limit int) ([]*EventPrices, int, int, error)
	GetByID(ctx context.Context, id int) (*EventPrices, error)
	GetByRestaurantID(ctx context.Context, restaurantID int) ([]*EventPrices, error)
	Update(ctx context.Context, event *EventPrices) error
	Delete(ctx context.Context, id int) error
}
type EventPrices struct {
	Id                int
	RestaurantId      int
	EventType         string
	TablePrice        float64
	WaiterPrice       float64
	MaxGuests         int
	TableSeats        int
	MaxWaiters        int
	AlcoholPermission bool
}
type CreateEventPrices struct {
	RestaurantId      int
	EventType         string
	TablePrice        float64
	WaiterPrice       float64
	MaxGuests         int
	TableSeats        int
	MaxWaiters        int
	AlcoholPermission bool
}
type UpdateEventPrices struct {
	RestaurantId      int
	EventType         string
	TablePrice        float64
	WaiterPrice       float64
	MaxGuests         int
	TableSeats        int
	MaxWaiters        int
	AlcoholPermission bool
}
