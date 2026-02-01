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

type cartRepo struct {
    db *gorm.DB
}

func NewCartRepo(db *gorm.DB) domain.CartRepository {
    return &cartRepo{db: db}
}

func (r *cartRepo) AddItem(ctx context.Context, item *domain.CartItem) error {
    var existingItem domain.CartItem


    err := r.db.WithContext(ctx).
        Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
        First(&existingItem).Error

    if err == nil {

        existingItem.Quantity += item.Quantity
        return r.db.WithContext(ctx).Save(&existingItem).Error
    } else if errors.Is(err, gorm.ErrRecordNotFound) {

        return r.db.WithContext(ctx).Create(item).Error
    }


    return err
}

func (r *cartRepo) ClearCart(ctx context.Context, userID uint) error {
    return r.db.WithContext(ctx).
        Unscoped().
        Where("user_id = ?", userID).
        Delete(&domain.CartItem{}).Error
}

func (r *cartRepo) RemoveItem(ctx context.Context, userID, productID uint) error {
    result := r.db.WithContext(ctx).
        Unscoped(). 
        Where("user_id = ? AND product_id = ?", userID, productID).
        Delete(&domain.CartItem{})

    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("item not found in cart")
    }
    return nil
}

func (r *cartRepo) GetCart(ctx context.Context, userID uint) ([]domain.CartItem, error) {
    var items []domain.CartItem
    
    err := r.db.WithContext(ctx).
        Preload("Product").
        Where("user_id = ?", userID).
        Find(&items).Error
        
    return items, err
}

func (r *cartRepo) UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error {

    result := r.db.WithContext(ctx).
        Model(&domain.CartItem{}).
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