package handler

import (
	"github.com/gofiber/fiber/v2"
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
        return response.Response(c,fiber.StatusUnauthorized,"no user found",nil,nil)
    }

    err := h.svc.PlaceOrder(c.Context(), uint(userIDFloat))
    if err != nil {
        return response.Response(c,fiber.StatusBadRequest,"server",nil,err.Error())
    }
	return response.Response(c,fiber.StatusOK,"Order placed successfully",nil,nil)
}