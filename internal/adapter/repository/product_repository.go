package repository

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

type Product struct {
	gorm.Model
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	Name        string         `json:"name" validate:"required"`
	Price       int            `json:"price" validate:"required"`
	Description string         `json:"desc" validate:"required"`
	Category    string         `json:"category" validate:"required"`
	Offer       string         `json:"offer,omitempty"`
	OfferPrice  int            `json:"offerprice,omitempty"`
	Production  string         `json:"production,omitempty"`
	Stock       uint           `json:"stock"`
}

func (p *Product) ToDomain() *domain.Product {
	return &domain.Product{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: func() *time.Time {
			if p.DeletedAt.Valid {
				return &p.DeletedAt.Time
			}
			return nil
		}(),
		Images:      []string(p.Images),
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Category:    p.Category,
		Offer:       p.Offer,
		OfferPrice:  p.OfferPrice,
		Production:  p.Production,
		Stock:       p.Stock,
	}
}

func fromDomainProduct(p *domain.Product) *Product {
	var deletedAt gorm.DeletedAt
	if p.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	}
	return &Product{
		Model: gorm.Model{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: deletedAt,
		},
		Images:      pq.StringArray(p.Images),
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		Category:    p.Category,
		Offer:       p.Offer,
		OfferPrice:  p.OfferPrice,
		Production:  p.Production,
		Stock:       p.Stock,
	}
}

func NewProductRepo(db *gorm.DB) domain.ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	dbProduct := fromDomainProduct(product)
	err := r.db.WithContext(ctx).Create(dbProduct).Error
	if err != nil {
		return err
	}
	// Update original product with ID and timestamps
	product.ID = dbProduct.ID
	product.CreatedAt = dbProduct.CreatedAt
	product.UpdatedAt = dbProduct.UpdatedAt
	return nil
}

func (r *productRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	var dbProducts []Product

	err := r.db.WithContext(ctx).Find(&dbProducts).Error
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}

	return products, nil
}

func (r *productRepo) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var dbProduct Product

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbProduct).Error
	if err != nil {
		return nil, err
	}
	return dbProduct.ToDomain(), err
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	dbProduct := fromDomainProduct(product)
	return r.db.WithContext(ctx).Save(dbProduct).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (r *productRepo) GetByProduction(ctx context.Context, status string) ([]domain.Product, error) {
	var dbProducts []Product

	err := r.db.WithContext(ctx).
		Where("Production = ?", status).
		Find(&dbProducts).Error

	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}

	return products, nil
}
