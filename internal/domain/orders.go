package domain

import (
	"context"
)

type OrderItem struct {
	OrderId uint
	Order   Order `gorm:"foreignkey:OrderId;references:ID"`

	Image string

	ProductId uint
	Product   Product `gorm:"foreignkey:ProductId;references:ID"`

	Quantity uint
	Price    float64
}

type Order struct {
	ID          uint
	UserID      uint   `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Status      string `json:"status"`
	Quantity    uint
	TotalAmount float64 `json:"total"`
}

type SalesData struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userid uint, Orders []OrderItem) error
	GetOrdersByUserIDAndOrderID(ctx context.Context, userID, OrderID uint) ([]Order, error)
	GetOrdersByOrderID(ctx context.Context, OrderID uint) ([]OrderItem, error)
	GetAllOrdersByUserID(ctx context.Context, userID uint) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
	CancelOrder(ctx context.Context, orderID, userID uint) error
	GetTotalSalesByDate(ctx context.Context) ([]SalesData, error)
	GetOrderCountsByStatus(ctx context.Context) ([]StatusCount, error)
}
