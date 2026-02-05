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
	BlockUser(ctx context.Context, userId uint) (string, error)
	AddNewProduct(ctx context.Context, newProduct *domain.Product) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	DeleteUser(ctx context.Context, userID uint) error
	Production(ctx context.Context, status string) ([]domain.Product, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	SearchProducts(ctx context.Context, query string) ([]domain.Product, error)
	SearchUsers(ctx context.Context, query string) ([]domain.User, error)
	FilterProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
}

type adminUseCase struct {
	repo        domain.UserRepository
	productrepo domain.ProductRepository
	oredersRepo domain.OrderRepository
}

func NewAdminUseCase(rep domain.UserRepository, productrepo domain.ProductRepository, oredersRepo domain.OrderRepository) AdminUseCase {
	return &adminUseCase{
		repo:        rep,
		productrepo: productrepo,
		oredersRepo: oredersRepo,
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
) (string, error) {

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return "", err
	}

	// Toggle block status
	user.IsBlocked = !user.IsBlocked

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