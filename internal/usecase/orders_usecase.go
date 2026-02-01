package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type OrderService interface {
    PlaceOrder(ctx context.Context, userID uint) error
}

type orderService struct {
    orderRepo domain.OrderRepository
    cartRepo  domain.CartRepository 
}

func NewOrderService(oRepo domain.OrderRepository, cRepo domain.CartRepository) OrderService {
    return &orderService{
        orderRepo: oRepo,
        cartRepo:  cRepo,
    }
}

func (s *orderService) PlaceOrder(ctx context.Context, userID uint) error {
    // 1. Get Cart Items
    cartItems, err := s.cartRepo.GetCart(ctx, userID)
    if err != nil {
        return err
    }
    if len(cartItems) == 0 {
        return errors.New("cart is empty")
    }

    // 2. Calculate Total & Prepare Order Items
    var totalAmount float64
    var orderItems []domain.OrderItem

    for _, cartItem := range cartItems {
        // Check Stock
        if cartItem.Product.Stock < cartItem.Quantity {
            return errors.New("product " + cartItem.Product.Name + " is out of stock")
        }

        totalAmount += float64(cartItem.Product.Price) * float64(cartItem.Quantity)

        orderItems = append(orderItems, domain.OrderItem{
            ProductID: cartItem.ProductID,
            Quantity:  int(cartItem.Quantity),
            Price:     float64(cartItem.Product.Price), 
        })
    }


    order := &domain.Order{
        UserID:      userID,
        TotalAmount: totalAmount,
        Status:      "Pending",
        Items:       orderItems,
    }

    return s.orderRepo.CreateOrder(ctx, order)
}