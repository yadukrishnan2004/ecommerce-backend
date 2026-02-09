package repository

import (
	"context"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Name      string `json:"name" validate:"required"` 
	Phone     string `json:"phone" validate:"required,min=10"`
	HouseName string `json:"house_name" validate:"required"` 
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	PinCode   string `json:"pin_code" validate:"required"`
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) domain.AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) Create(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *addressRepo) GetByUserID(ctx context.Context, userID uint) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}