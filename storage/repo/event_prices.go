package repo

import "context"

type EventPricesI interface {
	Create(ctx context.Context, req *EventPrice) (*EventPrice, error)
	GetAll(ctx context.Context, restaurantID, eventType string, page, limit int) ([]*EventPrice, int, int, error)
	GetByID(ctx context.Context, id int) (*EventPrice, error)
	GetByRestaurantID(ctx context.Context, restaurantID int) ([]*EventPrice, error)
	GetAllByRestaurantIDs(ctx context.Context, restaurantMap map[int]bool, eventType string, page, limit int) ([]EventPrice, int, int, error)
	Update(ctx context.Context, event *EventPrice) error
	Delete(ctx context.Context, id int) error
}

// type EventPrices struct {
// 	Id                int
// 	RestaurantId      int
// 	EventType         string
// 	TablePrice        float64
// 	WaiterPrice       float64
// 	MaxGuests         int
// 	TableSeats        int
// 	MaxWaiters        int
// 	AlcoholPermission bool
// }
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

// EventPrices jadvali
type EventPrice struct {
	Id                uint       `gorm:"primaryKey"`
	RestaurantId      uint       `gorm:"not null"`
	EventType         string     `gorm:"type:enum('morning','night');not null"`
	TablePrice        float64    `gorm:"type:decimal(12,2);not null"`
	WaiterPrice       float64    `gorm:"type:decimal(12,2);not null"`
	MaxGuests         int        `gorm:"not null"`
	TableSeats        int        `gorm:"not null"`
	MaxWaiters        int        `gorm:"not null"`
	AlcoholPermission bool       `gorm:"not null"`
	Restaurant        Restaurant `json:"-" gorm:"foreignKey:RestaurantId;constraint:OnDelete:CASCADE"`
}
