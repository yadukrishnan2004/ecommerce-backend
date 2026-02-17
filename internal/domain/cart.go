package domain

import (
	"context"

	"github.com/lib/pq"
)

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  uint    `json:"quantity"`
}

type CartItemView struct {
	Image       pq.StringArray `gorm:"type:text[]"`
	CartID      uint
	Quantity    uint
	ProductID   uint
	Name        string
	Price       int64
	Offer       string
	OfferPrice  int
	Description string
	Stock       uint
}

type CartRepository interface {
	AddItem(ctx context.Context, item *CartItem) error
	ClearCart(ctx context.Context, userID uint) error
	RemoveItem(ctx context.Context, userID, productID uint) error
	GetCart(ctx context.Context, userID uint) ([]CartItemView, error)
	UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error
}
