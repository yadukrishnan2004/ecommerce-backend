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
	Quantity    uint	`json:"quantity"`
	TotalAmount float64 `json:"total"`
}

type OrderItem struct {
	gorm.Model

	OrderId uint
	Order   Order `gorm:"foreignkey:OrderId;references:ID"`

    image   string

	ProductId uint
	Product   Product `gorm:"foreignkey:ProductId;references:ID"`

	Quantity uint
	Price    float64
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

	tx := r.db.WithContext(ctx).Begin()
	dbOrder := Order{
			UserID: userid,
			Status: "pending",
			Quantity: uint(len(Orders)),
			TotalAmount: totalPrice,
	}

	if err := tx.Create(&dbOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i:= range Orders{
		Orders[i].OrderId=dbOrder.ID
	}

	if err := tx.Create(&Orders).Error; err != nil {
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
	var dbOrders []domain.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id", userID, OrderID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}
	return dbOrders, nil
}

func (r *orderRepo) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	var dbOrders []domain.Order

	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	return dbOrders, nil
}

func (r *orderRepo) GetAllOrdersByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var dbOrders []domain.Order

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}
	return dbOrders, nil
}


func (r *orderRepo) GetOrdersByOrderID(ctx context.Context, OrderID uint) ([]domain.OrderItem, error) {
	var OrderItem []domain.OrderItem
	err:=r.db.WithContext(ctx).
	Where("OrderId = ?",OrderID).
	Order("created_at desc").
	Find(&OrderItem).Error

	if err != nil {
		return nil, err
	}
	return OrderItem, nil
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
		var dbOrder domain.Order

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

		items,err:=r.GetOrdersByOrderID(ctx,orderID)
		if err != nil {
			return err
		}

		for _, item := range items {
			if err := tx.Model(&domain.Product{}).
				Where("id = ?", item.ProductId).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
