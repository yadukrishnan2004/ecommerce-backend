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
	ID                uint
	UserID            uint    `json:"user_id"`
	User              User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	AddressID         uint    `json:"address_id"`
	Address           Address `json:"address" gorm:"foreignKey:AddressID;references:ID"`
	Status            string  `json:"status"`
	Quantity          uint
	TotalAmount       float64 `json:"total"`
	PaymentMethod     string  `json:"payment_method"`
	RazorpayOrderID   string  `json:"razorpay_order_id"`
	RazorpayPaymentID string  `json:"razorpay_payment_id"`
}

type SalesData struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TopProduct struct {
	ProductID    uint    `json:"product_id"`
	Name         string  `json:"name"`
	QuantitySold uint    `json:"quantity_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userid uint, addressID uint, Orders []OrderItem, paymentMethod, razorpayOrderID, razorpayPaymentID string) error
	GetOrdersByUserIDAndOrderID(ctx context.Context, userID, OrderID uint) ([]Order, error)
	GetOrdersByOrderID(ctx context.Context, OrderID uint) ([]OrderItem, error)
	GetOrderItemsByOrderIDs(ctx context.Context, orderIDs []uint) ([]OrderItem, error)
	GetAllOrdersByUserID(ctx context.Context, userID uint, limit, offset int) ([]Order, error)
	GetAllOrders(ctx context.Context, limit, offset int) ([]Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
	UpdateStatusByRazorpayOrderID(ctx context.Context, razorpayOrderID string, status string) error
	CancelOrder(ctx context.Context, orderID, userID uint) error
	GetTotalSalesByDate(ctx context.Context) ([]SalesData, error)
	GetOrderCountsByStatus(ctx context.Context) ([]StatusCount, error)
	GetDashboardMetrics(ctx context.Context) (totalRevenue float64, totalOrders int64, averageOrderValue float64, err error)
	GetTopSellingProducts(ctx context.Context, limit int) ([]TopProduct, error)
}
