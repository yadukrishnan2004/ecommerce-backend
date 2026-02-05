package domain

import "context"



type Wishlist struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}



type WishlistRepository interface {
	Add(ctx context.Context, item *Wishlist) error
	Remove(ctx context.Context, userID, productID uint) error
	DeleteAll(ctx context.Context, userID uint) error
	GetAll(ctx context.Context, userID uint) ([]Wishlist, error)
}