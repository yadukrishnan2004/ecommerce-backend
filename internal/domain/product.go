package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
    Images      []string   `json:"images"`
	Name        string     `json:"name" validate:"required"`
	Price       int        `json:"price" validate:"required"`
	Description string     `json:"desc" validate:"required"`
	Category    string     `json:"category" validate:"required"`
	Offer       string     `json:"offer,omitempty"`
	OfferPrice  int        `json:"offerprice,omitempty"`
	Production  string     `json:"production,omitempty"`
	Stock       uint       `json:"stock"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id uint) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	GetByProduction(ctx context.Context, name string) ([]Product, error)
}
