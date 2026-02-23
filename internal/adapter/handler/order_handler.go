package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type OrderHandler struct {
	svc usecase.OrderService
}

func NewOrderHandler(svc usecase.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusUnauthorized, "no user found", nil, nil)
	}

	type PlaceOrderReq struct {
		AddressID     uint   `json:"address_id"`
		PaymentMethod string `json:"payment_method"`
	}
	req := new(PlaceOrderReq)
	if err := c.BodyParser(req); err != nil {
		return response.Response(c, fiber.StatusBadRequest, "invalid input", nil, nil)
	}

	razorpayOrderID, err := h.svc.PlaceOrder(c.Context(), uint(userIDFloat), req.AddressID, req.PaymentMethod)
	if err != nil {
		return response.Response(c, fiber.StatusBadRequest, "server", nil, err.Error())
	}

	resData := map[string]string{
		"razorpay_order_id": razorpayOrderID,
	}

	return response.Response(c, fiber.StatusOK, "Order placed successfully", resData, nil)
}

func (h *OrderHandler) GetOrderHistory(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusUnauthorized, "no user found", nil, nil)
	}

	orders, err := h.svc.GetOrderHistory(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, fiber.StatusBadRequest, "server", nil, err.Error())
	}
	history := dto.Orders{
		Count: len(orders),
		Items: orders,
	}

	return response.Response(c, fiber.StatusOK, "get Order successfully", history, nil)

}

func (h *OrderHandler) BuyNow(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusUnauthorized, "no user found", nil, nil)
	}

	type BuyNowReq struct {
		ProductID     uint   `json:"product_id"`
		Quantity      int    `json:"quantity"`
		AddressID     uint   `json:"address_id"`
		PaymentMethod string `json:"payment_method"`
	}
	req := new(BuyNowReq)
	if err := c.BodyParser(req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input value", nil, nil)
	}

	if req.Quantity <= 0 {
		return response.Response(c, http.StatusBadRequest, "Quantity must be greater than 0", nil, nil)
	}

	razorpayOrderID, err := h.svc.BuyNow(c.Context(), uint(userIDFloat), req.AddressID, req.ProductID, req.Quantity, req.PaymentMethod)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "server error", nil, err.Error())
	}

	resData := map[string]string{
		"razorpay_order_id": razorpayOrderID,
	}

	return response.Response(c, http.StatusOK, "Order placed successfully", resData, nil)
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusUnauthorized, "no user found", nil, nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, fiber.StatusBadRequest, "invalid order id", nil, nil)
	}

	orderItems, err := h.svc.GetOrderDetails(c.Context(), uint(userIDFloat), uint(id))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to get order details", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Order details fetched successfully", orderItems, nil)
}

func (h *OrderHandler) VerifyPayment(c *fiber.Ctx) error {
	type VerifyPaymentReq struct {
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}
	req := new(VerifyPaymentReq)
	if err := c.BodyParser(req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input value", nil, nil)
	}

	err := h.svc.VerifyPayment(c.Context(), req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "payment verification failed", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Payment verified successfully", nil, nil)
}
