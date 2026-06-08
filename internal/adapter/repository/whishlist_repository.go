package repository

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Wishlist struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

func (w *Wishlist) ToDomain() *domain.Wishlist {
	return &domain.Wishlist{
		ID:        w.ID,
		UserID:    w.UserID,
		ProductID: w.ProductID,
		Product:   *w.Product.ToDomain(),
	}
}

func FromDomainWishlist(w *domain.Wishlist) *Wishlist {
	return &Wishlist{
		Model: gorm.Model{
			ID: w.ID,
		},
		UserID:    w.UserID,
		ProductID: w.ProductID,
	}
}

type wishlistRepo struct {
	db *gorm.DB
}

func NewWishlistRepo(db *gorm.DB) domain.WishlistRepository {
	return &wishlistRepo{db: db}
}

func (r *wishlistRepo) Add(ctx context.Context, item *domain.Wishlist) error {
	var existing Wishlist

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existing).Error

	if err == nil {
		return errors.New("item already in wishlist")
	}

	newItem := FromDomainWishlist(item)
	if err := r.db.WithContext(ctx).Create(newItem).Error; err != nil {
		return err
	}
	item.ID = newItem.ID
	return nil
}

func (r *wishlistRepo) Remove(ctx context.Context, userID, productID uint) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&Wishlist{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("wish list don't have an item")
	}
	return nil
}

func (r *wishlistRepo) DeleteAll(ctx context.Context, userID uint) error {
	//hard delete
	return r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ?", userID).
		Delete(&Wishlist{}).Error
}

func (r *wishlistRepo) GetAll(ctx context.Context, userID uint) ([]domain.WishlistItemView, error) {
	var items []domain.WishlistItemView

	err := r.db.WithContext(ctx).
		Table("wishlists").
		Select(`
			wishlists.id as wishlist_id, 
			products.id as product_id, 
			COALESCE(array_agg(product_images.image_url) FILTER (WHERE product_images.image_url IS NOT NULL), '{}') as images, 
			products.category, 
			products.name, 
			products.price, 
			products.offer, 
			products.offer_price, 
			products.description, 
			products.stock, 
			products.production
		`).
		Joins("JOIN products ON products.id = wishlists.product_id").
		Joins("LEFT JOIN product_images ON product_images.product_id = products.id").
		Where("wishlists.user_id = ? AND wishlists.deleted_at IS NULL", userID).
		Group("wishlists.id, products.id").
		Scan(&items).Error

	return items, err
}
