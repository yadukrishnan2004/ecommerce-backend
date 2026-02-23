package usecase

import (
	"context"
	"errors"

	"github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, userID uint, addressID uint, paymentMethod string) (string, error)
	GetOrderHistory(ctx context.Context, userID uint) ([]domain.Order, error)
	BuyNow(ctx context.Context, userID, addressID, productID uint, quantity int, paymentMethod string) (string, error)
	GetOrderDetails(ctx context.Context, userID, orderID uint) ([]domain.OrderItem, error)
	VerifyPayment(ctx context.Context, razorpayOrderID, razorpayPaymentID, razorpaySignature string) error
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

func (s *orderService) PlaceOrder(ctx context.Context, userID uint, addressID uint, paymentMethod string) (string, error) {
	cartItems, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return "", err
	}
	if len(cartItems) == 0 {
		return "", errors.New("cart is empty")
	}

	// var totalAmount float64
	var orderItems []domain.OrderItem

	for _, cartItem := range cartItems {
		if cartItem.Stock < cartItem.Quantity {
			return "", errors.New("product " + cartItem.Name + " is out of stock")
		}
		Pro, err := s.Product.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return "", errors.New("no paroduct found")
		}

		// totalAmount += float64(cartItem.Price) * float64(cartItem.Quantity)

		orderItems = append(orderItems, domain.OrderItem{
			ProductId: Pro.ID,
			Image:     Pro.Images[0],
			Quantity:  cartItem.Quantity,
			Price:     Pro.Price,
		})
	}
	var grandTotal float64
	var totalPrice float64
	for _, item := range orderItems {
		totalPrice += item.Price * float64(item.Quantity)
	}
	shipping := 50.0
	tax := totalPrice * 0.1
	grandTotal = totalPrice + shipping + tax

	var razorpayOrderID string
	if paymentMethod == "Razorpay" {
		client := razorpay.NewClient(config.Load().RAZORPAY_KEY, config.Load().RAZORPAY_SECRET)
		data := map[string]interface{}{
			"amount":   int(grandTotal * 100), // convert to paise
			"currency": "INR",
			"receipt":  "receipt_order_1", // In production use a random/unique receipt
		}
		body, err := client.Order.Create(data, nil)
		if err != nil {
			return "", errors.New("failed to generate razorpay order: " + err.Error())
		}
		razorpayOrderID = body["id"].(string)
	}

	err = s.orderRepo.CreateOrder(ctx, userID, addressID, orderItems, paymentMethod, razorpayOrderID, "")
	if err != nil {
		return "", err
	}
	return razorpayOrderID, nil
}

func (s *orderService) GetOrderHistory(ctx context.Context, userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetAllOrdersByUserID(ctx, userID)
}

func (s *orderService) BuyNow(ctx context.Context, userID, addressID, productID uint, quantity int, paymentMethod string) (string, error) {
	Pro, err := s.Product.GetByID(ctx, productID)
	if err != nil {
		return "", errors.New("product not found")
	}

	if Pro.Stock < uint(quantity) {
		return "", errors.New("insufficient stock")
	}
	var orderItems []domain.OrderItem
	orderItems = append(orderItems, domain.OrderItem{
		ProductId: Pro.ID,
		Image:     Pro.Images[0],
		Quantity:  uint(quantity),
		Price:     Pro.Price,
	})
	var grandTotal float64
	totalPrice := Pro.Price * float64(quantity)
	shipping := 50.0
	tax := totalPrice * 0.1
	grandTotal = totalPrice + shipping + tax

	var razorpayOrderID string
	if paymentMethod == "Razorpay" {
		client := razorpay.NewClient(config.Load().RAZORPAY_KEY, config.Load().RAZORPAY_SECRET)
		data := map[string]interface{}{
			"amount":   int(grandTotal * 100), // convert to paise
			"currency": "INR",
			"receipt":  "receipt_buynow_1", // In production use a unique receipt
		}
		body, err := client.Order.Create(data, nil)
		if err != nil {
			return "", errors.New("failed to generate razorpay order: " + err.Error())
		}
		razorpayOrderID = body["id"].(string)
	}

	err = s.orderRepo.CreateOrder(ctx, userID, addressID, orderItems, paymentMethod, razorpayOrderID, "")
	if err != nil {
		return "", err
	}
	return razorpayOrderID, nil
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

func (s *orderService) VerifyPayment(ctx context.Context, razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {
	secret := config.Load().RAZORPAY_SECRET
	params := map[string]interface{}{
		"razorpay_order_id":   razorpayOrderID,
		"razorpay_payment_id": razorpayPaymentID,
	}

	isValid := utils.VerifyPaymentSignature(params, razorpaySignature, secret)
	if !isValid {
		return errors.New("invalid payment signature")
	}

	// Update order status to paid in DB if necessary based on razorpay_order_id
	// (Needs new OrderRepo func to update by RazorpayOrderID, we'll keep it simple or implement next)
	err := s.orderRepo.UpdateStatusByRazorpayOrderID(ctx, razorpayOrderID, "Placed")
	if err != nil {
		return errors.New("failed to update order status")
	}

	return nil
}
