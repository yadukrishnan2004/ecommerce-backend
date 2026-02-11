package usecase

import (
	"context"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)


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

type AddressUsecase interface {
	AddAddress(ctx context.Context, address *domain.Address) error
	GetUserAddresses(ctx context.Context, userID uint) ([]domain.Address, error)
}


type addressUsecase struct {
	repo domain.AddressRepository
}

func NewAddressUsecase(repo domain.AddressRepository) AddressUsecase {
	return &addressUsecase{repo: repo}
}

func (s *addressUsecase) AddAddress(ctx context.Context, address *domain.Address) error {
	return s.repo.Create(ctx, address)
}

func (s *addressUsecase) GetUserAddresses(ctx context.Context, userID uint) ([]domain.Address, error) {
	return s.repo.GetByUserID(ctx, userID)
}