package domain

import (
	"context"
)


type CartItem struct {
    ID        uint    `json:"id" gorm:"primaryKey"`
    UserID    uint    `json:"user_id"`
    ProductID uint    `json:"product_id"`
    Product   Product `json:"product" gorm:"foreignKey:ProductID"` 
    Quantity  uint     `json:"quantity"`
}


type CartRepository interface {
    AddItem(ctx context.Context, item *CartItem) error
    ClearCart(ctx context.Context, userID uint) error
    RemoveItem(ctx context.Context, userID, productID uint) error
    GetCart(ctx context.Context, userID uint) ([]CartItem, error)
    UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error
}