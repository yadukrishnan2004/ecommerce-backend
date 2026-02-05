package domain

import (
	"context"
	"time"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
}


type OrderItem struct {
	gorm.Model
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	CreateSingleOrder(ctx context.Context, order *Order) error
	GetOrdersByUserID(ctx context.Context, userID uint) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
	GetByIDAndUser(ctx context.Context, orderID, userID uint) (*Order, error)
	CancelOrder(ctx context.Context, orderID, userID uint) error
}
