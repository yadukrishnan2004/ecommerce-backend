package dto

import "github.com/yadukrishnan2004/ecommerce-backend/internal/domain"

type Wishlist struct {
	Item  []domain.WishlistItemView `json:"item"`
	Count int                       `json:"count"`
}

type WishlistItemView struct {
	WishlistID uint `json:"wishlist_id"`
	ProductID  uint `json:"product_id"`

	Images   []string `json:"images"`
	Category string   `json:"category"`
	Name     string   `json:"name"`

	Price      int64  `json:"price"`
	Offer      string `json:"offer"`
	OfferPrice int64  `json:"offer_price"`

	Description string `json:"description"`
	Stock       uint   `json:"stock"`
	Production  string `json:"production"`
}
