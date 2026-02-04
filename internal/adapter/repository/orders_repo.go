package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type OrderProduct struct {
	gorm.Model
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	Name        string         `json:"name"`
	Price       int            `json:"price"`
	Description string         `json:"desc"`
	Category    string         `json:"category"`
	Offer       string         `json:"offer"`
	OfferPrice  int            `json:"offerprice"`
	Production  string         `json:"production"`
	Stock       uint           `json:"stock"`
}

func (OrderProduct) TableName() string {
	return "products"
}

func (p *OrderProduct) ToDomain() *domain.Product {
	return &domain.Product{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: func() *time.Time {
			if p.DeletedAt.Valid {
				return &p.DeletedAt.Time
			}
			return nil
		}(),
		Images:      []string(p.Images),
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Category:    p.Category,
		Offer:       p.Offer,
		OfferPrice:  p.OfferPrice,
		Production:  p.Production,
		Stock:       p.Stock,
	}
}

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
	ID        uint         `json:"id" gorm:"primaryKey"`
	OrderID   uint         `json:"order_id"`
	ProductID uint         `json:"product_id"`
	Product   OrderProduct `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int          `json:"quantity"`
	Price     float64      `json:"price"`
}

func (o *Order) ToDomain() *domain.Order {
	var domainItems []domain.OrderItem
	for _, i := range o.Items {
		domainItems = append(domainItems, domain.OrderItem{
			Model: gorm.Model{
				ID:        i.ID,
				CreatedAt: i.CreatedAt,
				UpdatedAt: i.UpdatedAt,
				DeletedAt: i.DeletedAt,
			},
			OrderID:   i.OrderID,
			ProductID: i.ProductID,
			Product:   *i.Product.ToDomain(),
			Quantity:  i.Quantity,
			Price:     i.Price,
		})
	}

	domainUser := *o.User.ToDomain()

	return &domain.Order{
		Model: gorm.Model{
			ID:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			DeletedAt: o.DeletedAt,
		},
		ID:          o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		Items:       domainItems,
		User:        domainUser,
	}
}

func fromDomainOrder(o *domain.Order) *Order {
	var items []OrderItem
	for _, i := range o.Items {
		items = append(items, OrderItem{
			Model:     i.Model,
			OrderID:   i.OrderID,
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
			Price:     i.Price,
		})
	}

	return &Order{
		Model:       o.Model,
		ID:          o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		Items:       items,
	}
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	dbOrder := fromDomainOrder(order)

	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Create(dbOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	order.ID = dbOrder.ID
	order.CreatedAt = dbOrder.CreatedAt

	for i := range dbOrder.Items {
		dbOrder.Items[i].OrderID = dbOrder.ID
	}

	for _, item := range dbOrder.Items {
		if err := tx.Model(&OrderProduct{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("user_id = ?", order.UserID).Delete(&CartItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepo) GetOrdersByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var dbOrders []Order

	err := r.db.WithContext(ctx).
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, *o.ToDomain())
	}

	return orders, nil
}

func (r *orderRepo) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	var dbOrders []Order

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items.Product").
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, *o.ToDomain())
	}

	return orders, nil
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

func (r *orderRepo) GetByIDAndUser(ctx context.Context, orderID, userID uint) (*domain.Order, error) {
	var dbOrder Order

	err := r.db.WithContext(ctx).
		Preload("Items.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&dbOrder).Error

	if err != nil {
		return nil, err
	}
	return dbOrder.ToDomain(), nil
}

func (r *orderRepo) CancelOrder(ctx context.Context, orderID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var dbOrder Order

		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Preload("Items").
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

		for _, item := range dbOrder.Items {
			if err := tx.Model(&OrderProduct{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepo) CreateSingleOrder(ctx context.Context, order *domain.Order) error {
	dbOrder := fromDomainOrder(order)
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dbOrder).Error; err != nil {
			return err
		}

		// Stock update
		for _, item := range dbOrder.Items {
			if err := tx.Model(&OrderProduct{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
