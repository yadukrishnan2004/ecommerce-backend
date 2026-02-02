package repository

import (
	"context"
	"errors"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
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

type orderRepo struct {
    db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) domain.OrderRepository {
    return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
    tx := r.db.WithContext(ctx).Begin()
    
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return err
    }

    for _, item := range order.Items {

        item.OrderID = order.ID
        
        if err := tx.Create(&item).Error; err != nil {
            tx.Rollback()
            return err
        }
        if err := tx.Model(&domain.Product{}).
            Where("id = ?", item.ProductID).
            Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    if err := tx.Where("user_id = ?", order.UserID).Delete(&domain.CartItem{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (r *orderRepo) GetOrdersByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
    var orders []domain.Order

    err := r.db.WithContext(ctx).
        Preload("Items.Product"). 
        Where("user_id = ?", userID).
        Order("created_at desc"). 
        Find(&orders).Error

    return orders, err
}

func (r *orderRepo) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
    var orders []domain.Order

    err := r.db.WithContext(ctx).
        Preload("User").          
        Preload("Items.Product").
        Order("created_at desc").
        Find(&orders).Error

    return orders, err
}

func (r *orderRepo) UpdateStatus(ctx context.Context, orderID uint, status string) error {
    result := r.db.WithContext(ctx).
        Model(&domain.Order{}).
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

func (r *orderRepo) GetByIDAndUser(ctx context.Context, orderID, userID uint) (*domain.Order, error) {
    var order domain.Order

    err := r.db.WithContext(ctx).
        Preload("Items.Product").
        Where("id = ? AND user_id = ?", orderID, userID).
        First(&order).Error

    return &order, err
}

func (r *orderRepo) CancelOrder(ctx context.Context, orderID, userID uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var order domain.Order
        
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            Preload("Items").
            Where("id = ? AND user_id = ?", orderID, userID).
            First(&order).Error; err != nil {
            return err 
        }

        if order.Status != "Pending" {
            return errors.New("order cannot be cancelled (already processed)")
        }

        if err := tx.Model(&order).Update("status", "Cancelled").Error; err != nil {
            return err
        }

        for _, item := range order.Items {

            if err := tx.Model(&domain.Product{}).
                Where("id = ?", item.ProductID).
                Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
                return err
            }
        }

        return nil
    })
}