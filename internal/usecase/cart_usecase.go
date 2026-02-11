package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type CartService interface {
	AddToCart(ctx context.Context, userID, productID uint, quantity uint) error
	ClearCart(ctx context.Context, userID uint) error
	RemoveItem(ctx context.Context, userID, productID uint) error
	GetCart(ctx context.Context, userID uint) ([]domain.CartItemView, error)
	UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error
}

type cartService struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
}

func NewCartService(cRepo domain.CartRepository, pRepo domain.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cRepo,
		productRepo: pRepo,
	}
}

func (s *cartService) AddToCart(ctx context.Context, userID, productID uint, quantity uint) error {

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	item := &domain.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return s.cartRepo.AddItem(ctx, item)
}

func (s *cartService) ClearCart(ctx context.Context, userID uint) error {
    return s.cartRepo.ClearCart(ctx, userID)
}

func (s *cartService) RemoveItem(ctx context.Context, userID, productID uint) error {
    return s.cartRepo.RemoveItem(ctx, userID, productID)
}

func (s *cartService) GetCart(ctx context.Context, userID uint) ([]domain.CartItemView, error) {
    return s.cartRepo.GetCart(ctx, userID)
}

func (s *cartService) UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error {

    if quantity <= 0 {
        return errors.New("quantity must be greater than 0")
    }

    product, err := s.productRepo.GetByID(ctx, productID)
    if err != nil {
        return errors.New("product not found")
    }
    
    if product.Stock < uint(quantity){
        return errors.New("insufficient stock available")
    }

    return s.cartRepo.UpdateQuantity(ctx, userID, productID, quantity)
}