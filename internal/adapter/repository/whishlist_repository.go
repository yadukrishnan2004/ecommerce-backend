package repository

import (
	"context"
	"errors"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Wishlist struct {
	gorm.Model
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

type wishlistRepo struct {
	db *gorm.DB
}

func NewWishlistRepo(db *gorm.DB) domain.WishlistRepository {
	return &wishlistRepo{db: db}
}

func (r *wishlistRepo) Add(ctx context.Context, item *domain.Wishlist) error {
	var existing domain.Wishlist

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existing).Error

	if err == nil {
		return errors.New("item already in wishlist")
	}

	return r.db.WithContext(ctx).Create(item).Error
}

func (r *wishlistRepo) Remove(ctx context.Context, userID, productID uint) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&domain.Wishlist{}).Error
}

func (r *wishlistRepo) DeleteAll(ctx context.Context, userID uint) error {
	//hard delete
    return r.db.WithContext(ctx).
		Unscoped().
        Where("user_id = ?", userID).
        Delete(&domain.Wishlist{}).Error
}

func (r *wishlistRepo) GetAll(ctx context.Context, userID uint) ([]domain.Wishlist, error) {
    var items []domain.Wishlist
    
    err := r.db.WithContext(ctx).
        Preload("Product"). 
        Where("user_id = ?", userID).
        Find(&items).Error
        
    return items, err
}