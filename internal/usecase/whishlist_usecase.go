package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type WishlistService interface {
	AddToWishlist(ctx context.Context, userID, productID uint) error
	RemoveFromWishlist(ctx context.Context, userID, productID uint) error
	ClearWishlist(ctx context.Context, userID uint) error
	GetWishlist(ctx context.Context, userID uint) ([] domain.Wishlist, error)
}

type wishlistService struct {
	wishRepo    domain.WishlistRepository
	productRepo domain.ProductRepository
}

func NewWishlistService(wRepo domain.WishlistRepository, pRepo domain.ProductRepository) WishlistService {
	return &wishlistService{
		wishRepo:    wRepo,
		productRepo: pRepo,
	}
}

func (s *wishlistService) AddToWishlist(ctx context.Context, userID, productID uint) error {

	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	item := &domain.Wishlist{
		UserID:    userID,
		ProductID: productID,
	}

	return s.wishRepo.Add(ctx, item)
}

func (s *wishlistService) RemoveFromWishlist(ctx context.Context, userID, productID uint) error {
	return s.wishRepo.Remove(ctx, userID, productID)
}

func (s *wishlistService) ClearWishlist(ctx context.Context, userID uint) error {
    return s.wishRepo.DeleteAll(ctx, userID)
}

func (s *wishlistService) GetWishlist(ctx context.Context, userID uint) ([]domain.Wishlist, error) {
    return s.wishRepo.GetAll(ctx, userID)
}