package dto

import "github.com/yadukrishnan2004/ecommerce-backend/internal/domain"

type Wishlist struct {
	Item  []domain.Wishlist `json:"item"`
	Count int               `json:"count"`
}
