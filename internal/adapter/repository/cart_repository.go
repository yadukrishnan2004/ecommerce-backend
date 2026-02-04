package repository

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
}

func (c *CartItem) ToDomain() *domain.CartItem {
	return &domain.CartItem{
		ID:        c.ID,
		UserID:    c.UserID,
		ProductID: c.ProductID,
		Product:   *c.Product.ToDomain(),
		Quantity:  uint(c.Quantity),
	}
}

func fromDomainCartItem(c *domain.CartItem) *CartItem {
	return &CartItem{
		Model: gorm.Model{
			ID: c.ID,
		},
		UserID:    c.UserID,
		ProductID: c.ProductID,
		Quantity:  int(c.Quantity),
	}
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) domain.CartRepository {
	return &cartRepo{db: db}
}

func (r *cartRepo) AddItem(ctx context.Context, item *domain.CartItem) error {
	var existingItem CartItem

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existingItem).Error

	if err == nil {
		existingItem.Quantity += int(item.Quantity)
		return r.db.WithContext(ctx).Save(&existingItem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := fromDomainCartItem(item)
		return r.db.WithContext(ctx).Create(newItem).Error
	}

	return err
}

func (r *cartRepo) ClearCart(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ?", userID).
		Delete(&CartItem{}).Error
}

func (r *cartRepo) RemoveItem(ctx context.Context, userID, productID uint) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&CartItem{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found in cart")
	}
	return nil
}

func (r *cartRepo) GetCart(ctx context.Context, userID uint) ([]domain.CartItem, error) {
	var dbItems []CartItem

	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&dbItems).Error

	if err != nil {
		return nil, err
	}

	var items []domain.CartItem
	for _, item := range dbItems {
		items = append(items, *item.ToDomain())
	}

	return items, nil
}

func (r *cartRepo) UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error {

	result := r.db.WithContext(ctx).
		Model(&CartItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found in cart")
	}

	return nil
}
