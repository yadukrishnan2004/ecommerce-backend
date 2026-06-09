package dto

import "github.com/yadukrishnan2004/ecommerce-backend/internal/domain"

type Orders struct {
	Count int            `json:"count"`
	Items []domain.Order `json:"items"`
}