package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, userID uint) error
	GetOrderHistory(ctx context.Context, userID uint) ([]domain.Order, error)
	BuyNow(ctx context.Context, userID, productID uint, quantity int) error
	GetOrderDetails(ctx context.Context, userID, orderID uint) ([]domain.OrderItem, error)
}

type orderService struct {
	orderRepo domain.OrderRepository
	cartRepo  domain.CartRepository
	Product   domain.ProductRepository
}

func NewOrderService(oRepo domain.OrderRepository, cRepo domain.CartRepository, pRepo domain.ProductRepository) OrderService {
	return &orderService{
		orderRepo: oRepo,
		cartRepo:  cRepo,
		Product:   pRepo,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, userID uint) error {
	cartItems, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}
	if len(cartItems) == 0 {
		return errors.New("cart is empty")
	}

	// var totalAmount float64
	var orderItems []domain.OrderItem

	for _, cartItem := range cartItems {
		if cartItem.Stock < cartItem.Quantity {
			return errors.New("product " + cartItem.Name + " is out of stock")
		}
		Pro, err := s.Product.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return errors.New("no paroduct found")
		}

		// totalAmount += float64(cartItem.Price) * float64(cartItem.Quantity)

		orderItems = append(orderItems, domain.OrderItem{
			ProductId: Pro.ID,
			Image:     Pro.Images[0],
			Quantity:  cartItem.Quantity,
			Price:     Pro.Price,
		})
	}

	return s.orderRepo.CreateOrder(ctx, userID, orderItems)
}

func (s *orderService) GetOrderHistory(ctx context.Context, userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetAllOrdersByUserID(ctx, userID)
}

func (s *orderService) BuyNow(ctx context.Context, userID, productID uint, quantity int) error {
	Pro, err := s.Product.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	if Pro.Stock < uint(quantity) {
		return errors.New("insufficient stock")
	}
	var orderItems []domain.OrderItem
	orderItems = append(orderItems, domain.OrderItem{
		ProductId: Pro.ID,
		Image:     Pro.Images[0],
		Quantity:  uint(quantity),
		Price:     Pro.Price,
	})
	return s.orderRepo.CreateOrder(ctx, userID, orderItems)
}

func (s *orderService) GetOrderDetails(ctx context.Context, userID, orderID uint) ([]domain.OrderItem, error) {
	// 1. Verify ownership (optional but good practice)
	// For now, we rely on the repo to fetch by order ID, but strict ownership check might need fetching order first
	// Assuming GetOrdersByOrderID is enough, but strictly we should check if Order.UserID == userID
	// Let's trust the repo fetch for now or we can add a check if needed.
	// Actually frontend calls this with just orderID, but from 'MyOrders' which implies ownership.
	// To be safe, we should probably fetch order first or modify repo to check user_id too.
	// But sticking to the plan:

	items, err := s.orderRepo.GetOrdersByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Security check: ensure the first item (if any) belongs to the order that belongs to the user
	if len(items) > 0 {
		if items[0].Order.UserID != userID {
			return nil, errors.New("unauthorized")
		}
	}

	return items, nil
}
