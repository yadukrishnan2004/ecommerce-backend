package repository

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint    `json:"user_id"`
	User        User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Status      string  `json:"status" gorm:"default:pending"`
	Quantity    uint    `json:"quantity"`
	TotalAmount float64 `json:"total"`
}

type OrderItem struct {
	gorm.Model

	OrderId uint
	Order   Order `gorm:"foreignkey:OrderId;references:ID"`

	Image string

	ProductId uint
	Product   Product `gorm:"foreignkey:ProductId;references:ID"`

	Quantity uint
	Price    float64
}

func (o *Order) ToDomain() domain.Order {
	return domain.Order{
		ID:          o.ID,
		UserID:      o.UserID,
		User:        *o.User.ToDomain(),
		Status:      o.Status,
		Quantity:    o.Quantity,
		TotalAmount: o.TotalAmount,
	}
}

func (oi *OrderItem) ToDomain() domain.OrderItem {
	return domain.OrderItem{
		OrderId:   oi.OrderId,
		Order:     oi.Order.ToDomain(),
		Image:     oi.Image,
		ProductId: oi.ProductId,
		Product:   *oi.Product.ToDomain(),
		Quantity:  oi.Quantity,
		Price:     oi.Price,
	}
}

func FromDomainOrderItem(oi domain.OrderItem) OrderItem {
	return OrderItem{
		OrderId:   oi.OrderId,
		Image:     oi.Image,
		ProductId: oi.ProductId,
		Quantity:  oi.Quantity,
		Price:     oi.Price,
	}
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(ctx context.Context, userid uint, Orders []domain.OrderItem) error {

	var totalPrice float64
	for _, Orderss := range Orders {
		totalPrice += float64(Orderss.Price)
	}

	shipping := 50.0
	tax := totalPrice * 0.1
	grandTotal := totalPrice + shipping + tax

	tx := r.db.WithContext(ctx).Begin()
	dbOrder := Order{
		UserID:      userid,
		Status:      "pending",
		Quantity:    uint(len(Orders)),
		TotalAmount: grandTotal,
	}

	if err := tx.Create(&dbOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	var dbOrderItems []OrderItem
	for _, item := range Orders {
		dbItem := FromDomainOrderItem(item)
		dbItem.OrderId = dbOrder.ID
		dbOrderItems = append(dbOrderItems, dbItem)
	}

	if err := tx.Create(&dbOrderItems).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", dbOrder.UserID).Delete(&CartItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepo) GetOrdersByUserIDAndOrderID(ctx context.Context, userID, OrderID uint) ([]domain.Order, error) {
	var dbOrders []Order
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, OrderID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}
	return orders, nil
}

func (r *orderRepo) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	var dbOrders []Order

	err := r.db.WithContext(ctx).
		Preload("User").
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}

	return orders, nil
}

func (r *orderRepo) GetAllOrdersByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var dbOrders []Order

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}
	return orders, nil
}

func (r *orderRepo) GetOrdersByOrderID(ctx context.Context, OrderID uint) ([]domain.OrderItem, error) {
	var dbOrderItems []OrderItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Order").
		Where("order_id = ?", OrderID).
		Order("created_at desc").
		Find(&dbOrderItems).Error

	if err != nil {
		return nil, err
	}

	var orderItems []domain.OrderItem
	for _, item := range dbOrderItems {
		orderItems = append(orderItems, item.ToDomain())
	}

	return orderItems, nil
}

func (r *orderRepo) UpdateStatus(ctx context.Context, orderID uint, status string) error {
	result := r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ?", orderID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *orderRepo) CancelOrder(ctx context.Context, orderID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var dbOrder Order

		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&dbOrder).Error; err != nil {
			return err
		}

		if dbOrder.Status != "Pending" {
			return errors.New("order cannot be cancelled (already processed)")
		}

		if err := tx.Model(&dbOrder).Update("status", "Cancelled").Error; err != nil {
			return err
		}

		items, err := r.GetOrdersByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		for _, item := range items {
			if err := tx.Model(&Product{}).
				Where("id = ?", item.ProductId).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepo) GetTotalSalesByDate(ctx context.Context) ([]domain.SalesData, error) {
	var sales []domain.SalesData
	// PostgreSQL: TO_CHAR(created_at, 'YYYY-MM-DD')
	err := r.db.WithContext(ctx).
		Model(&Order{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, SUM(total_amount) as total").
		Where("status = ?", "Delivered"). // Only count delivered/completed orders
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date asc").
		Limit(30). // Last 30 days
		Scan(&sales).Error

	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *orderRepo) GetOrderCountsByStatus(ctx context.Context) ([]domain.StatusCount, error) {
	var counts []domain.StatusCount
	err := r.db.WithContext(ctx).
		Model(&Order{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}
	return counts, nil
}
