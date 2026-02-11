package usecase

import (
	"context"
	"errors"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type OrderService interface {
    PlaceOrder(ctx context.Context, userID uint) error
    GetOrderHistory(ctx context.Context, userID uint) ([] domain.Order, error)  
    BuyNow(ctx context.Context, userID, productID uint, quantity int) error 
}

type orderService struct {
    orderRepo domain.OrderRepository
    cartRepo  domain.CartRepository 
    Product   domain.ProductRepository
}

func NewOrderService(oRepo domain.OrderRepository, cRepo domain.CartRepository,pRepo  domain.ProductRepository) OrderService {
    return &orderService{
        orderRepo: oRepo,
        cartRepo:  cRepo,
        Product: pRepo,
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
        Pro,err:=s.Product.GetByID(ctx,cartItem.ProductID)
        if err != nil {
            return errors.New("no paroduct found")
        }

        // totalAmount += float64(cartItem.Price) * float64(cartItem.Quantity)

        orderItems = append(orderItems, domain.OrderItem{
                    ProductId:Pro.ID ,
                    Image: Pro.Images[0],
                    Quantity: cartItem.Quantity,             
        })
    }

    return s.orderRepo.CreateOrder(ctx,userID,orderItems)
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
                    ProductId:Pro.ID ,
                    Image: Pro.Images[0],
                    Quantity: uint(quantity),
        })
    return s.orderRepo.CreateOrder(ctx,userID,orderItems)
}


