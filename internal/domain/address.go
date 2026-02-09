package domain

import "context"

type Address struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone" validate:"required,min=10"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	PinCode   string `json:"pin_code" validate:"required"`
}

type AddressRepository interface {
	Create(ctx context.Context, address *Address) error
	GetByUserID(ctx context.Context, userID uint) ([]Address, error)
}