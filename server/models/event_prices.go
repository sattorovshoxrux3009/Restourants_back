package models

type EventPrice struct {
	Id                int     `json:"id"`
	RestaurantId      int     `json:"restaurant_id"`
	EventType         string  `json:"event_type"`
	TablePrice        float64 `json:"table_price"`
	WaiterPrice       float64 `json:"waiter_price"`
	MaxGuests         int     `json:"max_guests"`
	TableSeats        int     `json:"table_seats"`
	MaxWaiters        int     `json:"max_waiters"`
	AlcoholPermission bool    `json:"alcohol_permission"`
}
type CreateEventPrices struct {
	RestaurantId      int     `json:"restaurant_id"`
	EventType         string  `json:"event_type"`
	TablePrice        float64 `json:"table_price"`
	WaiterPrice       float64 `json:"waiter_price"`
	MaxGuests         int     `json:"max_guests"`
	TableSeats        int     `json:"table_seats"`
	MaxWaiters        int     `json:"max_waiters"`
	AlcoholPermission bool    `json:"alcohol_permission"`
}
type UpdateEventPrices struct {
	TablePrice        float64 `json:"table_price"`
	WaiterPrice       float64 `json:"waiter_price"`
	MaxGuests         int     `json:"max_guests"`
	TableSeats        int     `json:"table_seats"`
	MaxWaiters        int     `json:"max_waiters"`
	AlcoholPermission bool    `json:"alcohol_permission"`
}
