package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type OrderService interface {
    PlaceOrder(ctx context.Context, userID uint) error
    GetOrderHistory(ctx context.Context, userID uint) ([] domain.Order, error)   
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
    cartItems, err := s.cartRepo.GetCart(ctx, userID)
    if err != nil {
        return err
    }
    if len(cartItems) == 0 {
        return errors.New("cart is empty")
    }

    var totalAmount float64
    var orderItems []domain.OrderItem

    for _, cartItem := range cartItems {
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

func (s *orderService) GetOrderHistory(ctx context.Context, userID uint) ([]domain.Order, error) {
    return s.orderRepo.GetOrdersByUserID(ctx, userID)
}


