package repository

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID            uint           `json:"user_id"`
	User              User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
	AddressID         uint           `json:"address_id"`
	Address           Address        `json:"address" gorm:"foreignKey:AddressID;references:ID"`
	Status            string         `json:"status" gorm:"default:pending"`
	TotalAmount       float64        `json:"total"`
	PaymentMethod     string         `json:"payment_method"`
	RazorpayOrderID   string         `json:"razorpay_order_id"`
	RazorpayPaymentID string         `json:"razorpay_payment_id"`
	Items             []OrderItem    `gorm:"foreignKey:OrderId"`
}

type OrderItem struct {
	gorm.Model

	OrderId uint
	Order   Order `gorm:"foreignkey:OrderId;references:ID"`

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

func (r *orderRepo) CreateOrder(ctx context.Context, userid uint, addressID uint, Orders []domain.OrderItem, paymentMethod, razorpayOrderID, razorpayPaymentID string) error {
	var totalPrice float64
	for _, Orderss := range Orders {
		totalPrice += float64(Orderss.Price)
	}

	shipping := 50.0
	tax := totalPrice * 0.1
	grandTotal := totalPrice + shipping + tax

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := "Pending"
	if paymentMethod == "COD" {
		status = "Placed"
	}

	dbOrder := Order{
		UserID:            userid,
		AddressID:         addressID,
		Status:            status,
		TotalAmount:       grandTotal,
		PaymentMethod:     paymentMethod,
		RazorpayOrderID:   razorpayOrderID,
		RazorpayPaymentID: razorpayPaymentID,
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

	for _, item := range Orders {
		var product Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", item.ProductId).First(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
		if product.Stock < item.Quantity {
			tx.Rollback()
			return errors.New("insufficient stock for product " + product.Name)
		}
		if err := tx.Model(&Product{}).Where("id = ?", item.ProductId).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("user_id = ?", dbOrder.UserID).Delete(&CartItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

		if dbOrder.Status != "Pending" && dbOrder.Status != "Placed" {
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

func (r *orderRepo) UpdateStatusByRazorpayOrderID(ctx context.Context, razorpayOrderID string, status string) error {
	return r.db.WithContext(ctx).Model(&Order{}).Where("razorpay_order_id = ?", razorpayOrderID).Update("status", status).Error
}
