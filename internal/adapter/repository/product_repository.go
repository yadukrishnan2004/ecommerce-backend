package repository

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

type Product struct{
	gorm.Model
    Name         string   `json:"name" validate:"required"`
    Price        int      `json:"price" validate:"required"`
    Description  string   `json:"desc" validate:"required"`
    Catogery     string   `json:"catoger" validate:"required"`
    Offer        string   `json:"offer,omitempty"`
    OfferPrice   int      `json:"offerprice,omitempty"`
	Production   string	  `json:"production,omitempty"`
	Stock		 uint	  `json:"stock"`
}

func NewProductRepo(db *gorm.DB) domain.ProductRepository {
	return &productRepo{db: db}
}


func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	
	err := r.db.WithContext(ctx).Find(&products).Error
		
	return products, err
}


func (r *productRepo) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product

	err := r.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepo) GetByProduction(ctx context.Context, status string) ([]domain.Product, error) {
    var product []domain.Product

    err := r.db.WithContext(ctx).
        Where("Production = ?", status).
        Find(&product).Error
        
    if err != nil {
        return nil, err
    }

    return product, nil
}


