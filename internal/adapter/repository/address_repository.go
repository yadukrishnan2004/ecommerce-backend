package repository

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

type PostalCode struct {
	PinCode string `gorm:"primaryKey;type:varchar(20)" json:"pin_code"`
	City    string `gorm:"type:varchar(100);not null" json:"city"`
	State   string `gorm:"type:varchar(100);not null" json:"state"`
}

type Address struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	Name       string     `json:"name" validate:"required"` 
	Phone      string     `json:"phone" validate:"required,min=10"`
	HouseName  string     `json:"house_name" validate:"required"` 
	Street     string     `json:"street" validate:"required"`
	PinCode    string     `json:"pin_code" validate:"required"`
	PostalCode PostalCode `gorm:"foreignKey:PinCode;references:PinCode"`
}

func (a *Address) ToDomain() *domain.Address {
	return &domain.Address{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Phone:     a.Phone,
		HouseName: a.HouseName,
		Street:    a.Street,
		City:      a.PostalCode.City,
		State:     a.PostalCode.State,
		PinCode:   a.PinCode,
	}
}

func fromDomainAddress(a *domain.Address) *Address {
	return &Address{
		Model: gorm.Model{
			ID: a.ID,
		},
		UserID:    a.UserID,
		Name:      a.Name,
		Phone:     a.Phone,
		HouseName: a.HouseName,
		Street:    a.Street,
		PinCode:   a.PinCode,
	}
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) domain.AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) Create(ctx context.Context, address *domain.Address) error {
	postalCode := PostalCode{
		PinCode: address.PinCode,
		City:    address.City,
		State:   address.State,
	}
	if err := r.db.WithContext(ctx).Save(&postalCode).Error; err != nil {
		return err
	}

	dbAddress := fromDomainAddress(address)
	err := r.db.WithContext(ctx).Create(dbAddress).Error
	if err != nil {
		return err
	}
	address.ID = dbAddress.ID
	return nil
}

func (r *addressRepo) GetByUserID(ctx context.Context, userID uint) ([]domain.Address, error) {
	var dbAddresses []Address
	err := r.db.WithContext(ctx).Preload("PostalCode").Where("user_id = ?", userID).Find(&dbAddresses).Error
	if err != nil {
		return nil, err
	}

	var addresses []domain.Address
	for _, a := range dbAddresses {
		addresses = append(addresses, *a.ToDomain())
	}
	return addresses, nil
}