package domain

import (
	"context"
	"time"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
    Images      []string   `json:"images" gorm:"serializer:json"`
	Name        string     `json:"name" validate:"required"`
	Price       float64        `json:"price" validate:"required"`
	Description string     `json:"desc" validate:"required"`
	Category    string     `json:"category" validate:"required"`
	Offer       string     `json:"offer,omitempty"`
	OfferPrice  int        `json:"offerprice,omitempty"`
	Production  string     `json:"production,omitempty"`
	Stock       uint       `json:"stock"`
}


type ProductFilter struct {
    Search    string  `json:"search"`
    MinPrice  float64 `json:"minprice"`
    MaxPrice  float64 `json:"maxprice"`
    Sort      string  `json:"sort"`
	Category  string  `json:"category"`
}


type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id uint) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	GetByProduction(ctx context.Context, name string) ([]Product, error)
	Search(ctx context.Context, query string) ([]Product, error)
	GetProducts(ctx context.Context, filter ProductFilter) ([]Product, error)
}
