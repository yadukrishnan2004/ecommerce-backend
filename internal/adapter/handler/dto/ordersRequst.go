package dto

import "github.com/yadukrishnan2004/ecommerce-backend/internal/domain"

type Orders struct {
	Count int
	Items []domain.Order
}