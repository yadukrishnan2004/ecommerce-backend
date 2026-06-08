package domain

import "context"

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAll(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id uint) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uint) error
	GetByName(ctx context.Context, name string) (*Category, error)
}
