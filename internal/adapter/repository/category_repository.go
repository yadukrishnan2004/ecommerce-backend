package repository

import (
	"context"
	"strings"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;type:varchar(100);not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
}

func (c *Category) ToDomain() *domain.Category {
	return &domain.Category{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
	}
}

func fromDomainCategory(c *domain.Category) *Category {
	return &Category{
		Model: gorm.Model{
			ID: c.ID,
		},
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
	}
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{db: db}
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) error {
	dbCategory := fromDomainCategory(category)
	dbCategory.Slug = slugify(category.Name)

	if err := r.db.WithContext(ctx).Create(dbCategory).Error; err != nil {
		return err
	}
	category.ID = dbCategory.ID
	category.Slug = dbCategory.Slug
	return nil
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]domain.Category, error) {
	var dbCategories []Category
	if err := r.db.WithContext(ctx).Order("name asc").Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	var categories []domain.Category
	for _, c := range dbCategories {
		categories = append(categories, *c.ToDomain())
	}
	return categories, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	var c Category
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return c.ToDomain(), nil
}

func (r *categoryRepo) Update(ctx context.Context, category *domain.Category) error {
	dbCategory := fromDomainCategory(category)
	dbCategory.Slug = slugify(category.Name)

	// Fetch existing first to retain timestamps if any or update GORM Model
	return r.db.WithContext(ctx).Model(&Category{}).Where("id = ?", category.ID).Updates(map[string]interface{}{
		"name":        dbCategory.Name,
		"slug":        dbCategory.Slug,
		"description": dbCategory.Description,
	}).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Category{}, id).Error
}

func (r *categoryRepo) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	var c Category
	if err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", strings.TrimSpace(name)).First(&c).Error; err != nil {
		return nil, err
	}
	return c.ToDomain(), nil
}
