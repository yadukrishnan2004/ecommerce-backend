package domain

import (
	"context"

	"github.com/lib/pq"
)



type Wishlist struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}


type WishlistItemView struct {
	WishlistID uint     `json:"wishlist_id"`
	ProductID  uint     `json:"product_id"`

	Images     pq.StringArray `gorm:"type:text[]"`
	Category   string   `json:"category"`
	Name       string   `json:"name"`

	Price      int64    `json:"price"`      
	Offer      string   `json:"offer"`
	OfferPrice int64    `json:"offer_price"` 

	Description string  `json:"description"`
	Stock       uint    `json:"stock"`
	Production  string  `json:"production"`
}



type WishlistRepository interface {
	Add(ctx context.Context, item *Wishlist) error
	Remove(ctx context.Context, userID, productID uint) error
	DeleteAll(ctx context.Context, userID uint) error
	GetAll(ctx context.Context, userID uint) ([]WishlistItemView, error)
}