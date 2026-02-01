package domain

import (
    "context"
    "time"
)

type Order struct {
    ID          uint        `json:"id" gorm:"primaryKey"`
    UserID      uint        `json:"user_id"`
    TotalAmount float64     `json:"total_amount"`
    Status      string      `json:"status"`
    CreatedAt   time.Time   `json:"created_at"`
    Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
    ID        uint    `json:"id" gorm:"primaryKey"`
    OrderID   uint    `json:"order_id"`
    ProductID uint    `json:"product_id"`
    Product   Product `json:"product" gorm:"foreignKey:ProductID"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"` 
}

type OrderRepository interface {
    CreateOrder(ctx context.Context, order *Order) error
}