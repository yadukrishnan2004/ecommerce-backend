package repository

import (
	"context"
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