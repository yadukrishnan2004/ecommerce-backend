package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

type ProductImage struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

type Product struct {
	gorm.Model
	Images      []ProductImage `json:"images" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Name        string         `json:"name" validate:"required"`
	Price       float64        `json:"price" validate:"required"`
	Description string         `json:"desc" validate:"required"`
	Category    string         `json:"category" validate:"required"`
	Offer       string         `json:"offer,omitempty"`
	OfferPrice  int            `json:"offerprice,omitempty"`
	Production  string         `json:"production,omitempty"`
	Stock       uint           `json:"stock"`
}

func (p *Product) ToDomain() *domain.Product {
	var images []string
	for _, img := range p.Images {
		images = append(images, img.ImageURL)
	}

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
		Images:      images,
		Name:        p.Name,
		Price:       float64(p.Price),
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

	var images []ProductImage
	for _, url := range p.Images {
		images = append(images, ProductImage{ImageURL: url})
	}

	return &Product{
		Model: gorm.Model{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: deletedAt,
		},
		Images:      images,
		Name:        p.Name,
		Price:       float64(p.Price),
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
	var catCount int64
	if err := r.db.WithContext(ctx).Table("categories").Where("LOWER(name) = LOWER(?)", strings.TrimSpace(product.Category)).Count(&catCount).Error; err != nil {
		return err
	}
	if catCount == 0 {
		return errors.New("category does not exist: " + product.Category)
	}

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

func (r *productRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	var dbProducts []Product

	query := r.db.WithContext(ctx).Preload("Images")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&dbProducts).Error
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

	err := r.db.WithContext(ctx).Preload("Images").Where("id = ?", id).First(&dbProduct).Error
	if err != nil {
		return nil, err
	}
	return dbProduct.ToDomain(), err
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	dbProduct := fromDomainProduct(product)
	// Delete existing images to avoid duplication
	if err := r.db.WithContext(ctx).Where("product_id = ?", dbProduct.ID).Delete(&ProductImage{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(dbProduct).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (r *productRepo) GetByProduction(ctx context.Context, status string, limit, offset int) ([]domain.Product, error) {
	var dbProducts []Product

	query := r.db.WithContext(ctx).Preload("Images").Where("Production = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&dbProducts).Error

	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}

	return products, nil
}

func (r *productRepo) Search(ctx context.Context, query string, limit, offset int) ([]domain.Product, error) {
	var dbProducts []Product

	searchPattern := "%" + query + "%"

	dbQuery := r.db.WithContext(ctx).Preload("Images").Where("name ILIKE ?", searchPattern)
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	err := dbQuery.Find(&dbProducts).Error

	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}

	return products, err
}

func (r *productRepo) GetProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	var dbProducts []Product

	query := r.db.WithContext(ctx).Preload("Images").Model(&Product{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	if filter.Category != "" {
		query = query.Where("category ILIKE ?", filter.Category)
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	switch filter.Sort {
	case "price_asc":
		query = query.Order("price asc")
	case "price_desc":
		query = query.Order("price desc")
	case "newest":
		query = query.Order("created_at desc")
	default:
		query = query.Order("id desc")
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&dbProducts).Error
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}

	return products, nil
}

func (r *productRepo) GetLowStockProducts(ctx context.Context, threshold int) ([]domain.Product, error) {
	var dbProducts []Product
	err := r.db.WithContext(ctx).Preload("Images").Where("stock <= ?", threshold).Order("stock asc").Find(&dbProducts).Error
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	for _, p := range dbProducts {
		products = append(products, *p.ToDomain())
	}
	return products, nil
}
