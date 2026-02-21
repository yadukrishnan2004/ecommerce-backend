package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type AdminUpdateUserInput struct {
	Name      *string
	Email     *string
	Role      *string
	IsActive  *bool
	IsBlocked *bool
}

type Product struct {
	Name        string `json:"name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Description string `json:"desc" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Offer       string `json:"offer,omitempty"`
	OfferPrice  int    `json:"offerprice,omitempty"`
}

type AdminUseCase interface {
	UpdateUser(ctx context.Context, userId uint, input AdminUpdateUserInput) (*domain.User, error)
	BlockUser(ctx context.Context, userId uint, blockedOpt *bool) (string, error)
	AddNewProduct(ctx context.Context, newProduct *domain.Product) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id uint, req *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	DeleteUser(ctx context.Context, userID uint) error
	Production(ctx context.Context, status string) ([]domain.Product, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	SearchProducts(ctx context.Context, query string) ([]domain.Product, error)
	SearchUsers(ctx context.Context, query string) ([]domain.User, error)
	FilterProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetDashboardGraphs(ctx context.Context) (map[string]interface{}, error)
	GetUserByID(ctx context.Context, userID uint) (*domain.User, error)
	GetUserCart(ctx context.Context, userID uint) ([]domain.CartItemView, error)
	GetUserWishlist(ctx context.Context, userID uint) ([]domain.WishlistItemView, error)
	GetUserAddresses(ctx context.Context, userID uint) ([]domain.Address, error)
}

// ... existing struct and constructor ...

// ... existing methods ...

func (s *adminUseCase) GetDashboardGraphs(ctx context.Context) (map[string]interface{}, error) {
	salesData, err := s.oredersRepo.GetTotalSalesByDate(ctx)
	if err != nil {
		return nil, err
	}

	statusCounts, err := s.oredersRepo.GetOrderCountsByStatus(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sales":  salesData,
		"orders": statusCounts,
	}, nil
}

type adminUseCase struct {
	repo        domain.UserRepository
	productrepo domain.ProductRepository
	oredersRepo domain.OrderRepository
	cartRepo    domain.CartRepository
	wishRepo    domain.WishlistRepository
	addressRepo domain.AddressRepository
}

func NewAdminUseCase(
	rep domain.UserRepository,
	productrepo domain.ProductRepository,
	oredersRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	wishRepo domain.WishlistRepository,
	addressRepo domain.AddressRepository,
) AdminUseCase {
	return &adminUseCase{
		repo:        rep,
		productrepo: productrepo,
		oredersRepo: oredersRepo,
		cartRepo:    cartRepo,
		wishRepo:    wishRepo,
		addressRepo: addressRepo,
	}
}

func (s *adminUseCase) UpdateUser(
	ctx context.Context,
	userId uint,
	input AdminUpdateUserInput,
) (*domain.User, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if input.IsBlocked != nil {
		user.IsBlocked = *input.IsBlocked
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *adminUseCase) BlockUser(
	ctx context.Context,
	userId uint,
	blockedOpt *bool,
) (string, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return "", err
	}

	// If an explicit blocked value is provided, honor it; otherwise toggle.
	if blockedOpt != nil {
		user.IsBlocked = *blockedOpt
	} else {
		user.IsBlocked = !user.IsBlocked
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return "", err
	}

	if user.IsBlocked {
		return "user blocked", nil
	}

	return "user unblocked", nil
}

func (s *adminUseCase) AddNewProduct(
	ctx context.Context,
	newProduct *domain.Product,
) error {
	if err := s.productrepo.Create(ctx, newProduct); err != nil {
		return errors.New("repository error on creating a repo")
	}
	return nil
}

func (s *adminUseCase) GetAllProducts(
	ctx context.Context,
) ([]domain.Product, error) {

	product, err := s.productrepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *adminUseCase) GetProduct(
	ctx context.Context,
	id uint,
) (*domain.Product, error) {
	product, err := s.productrepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *adminUseCase) UpdateProduct(
	ctx context.Context,
	id uint,
	req *domain.Product,
) error {
	// First, fetch the existing product
	product, err := s.productrepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update the allowed fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Offer != "" {
		product.Offer = req.Offer
	}
	if req.OfferPrice != 0 {
		product.OfferPrice = req.OfferPrice
	}
	if req.Production != "" {
		product.Production = req.Production
	}
	product.Stock = req.Stock // Assuming 0 is a valid stock to overwrite with

	// Usually images can optionally be overlaid/appended or replaced completely
	if len(req.Images) > 0 {
		product.Images = req.Images
	}

	if err := s.productrepo.Update(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *adminUseCase) DeleteProduct(
	ctx context.Context,
	id uint,
) error {
	err := s.productrepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *adminUseCase) DeleteUser(
	ctx context.Context,
	userID uint,
) error {
	return s.repo.Delete(ctx, userID)
}

func (s *adminUseCase) Production(
	ctx context.Context,
	status string,
) ([]domain.Product, error) {
	product, err := s.productrepo.GetByProduction(ctx, status)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *adminUseCase) UpdateStatus(ctx context.Context, id uint, status string) error {
	validStatuses := map[string]bool{
		"active":       true,
		"coming_soon":  true,
		"out_of_stock": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status value")
	}
	product, err := s.productrepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	product.Production = status
	return s.productrepo.Update(ctx, product)
}

func (s *adminUseCase) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	return s.oredersRepo.GetAllOrders(ctx)
}

func (s *adminUseCase) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {

	validStatuses := map[string]bool{
		"Pending":   true,
		"Shipped":   true,
		"Delivered": true,
		"Cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status value")
	}

	return s.oredersRepo.UpdateStatus(ctx, orderID, status)
}

func (s *adminUseCase) SearchProducts(ctx context.Context, query string) ([]domain.Product, error) {
	cleanQuery := strings.TrimSpace(query)

	if cleanQuery == "" {
		return []domain.Product{}, nil
	}
	return s.productrepo.Search(ctx, cleanQuery)
}

func (s *adminUseCase) SearchUsers(ctx context.Context, query string) ([]domain.User, error) {
	return s.repo.SearchUsers(ctx, query)
}

func (s *adminUseCase) FilterProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	return s.productrepo.GetProducts(ctx, filter)
}

func (s *adminUseCase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *adminUseCase) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *adminUseCase) GetUserCart(ctx context.Context, userID uint) ([]domain.CartItemView, error) {
	return s.cartRepo.GetCart(ctx, userID)
}

func (s *adminUseCase) GetUserWishlist(ctx context.Context, userID uint) ([]domain.WishlistItemView, error) {
	return s.wishRepo.GetAll(ctx, userID)
}

func (s *adminUseCase) GetUserAddresses(ctx context.Context, userID uint) ([]domain.Address, error) {
	return s.addressRepo.GetByUserID(ctx, userID)
}
