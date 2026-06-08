package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *UserHandler) GetOrder(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, nil)
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid user id", nil, nil)
	}

	order, err := h.svc.GetOrderDetail(c.Context(), uint(orderID), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "user not found", nil, nil)
	}

	return response.Response(c, http.StatusOK, "get order", order, nil)
}

func (h *UserHandler) GetOrderProduct(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, nil)
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid order id", nil, nil)
	}

	order, err := h.svc.GetOrderItemDetails(c.Context(), uint(orderID), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "order not found", nil, nil)
	}

	return response.Response(c, http.StatusOK, "get order", order, nil)
}

func (h *UserHandler) CancelOrder(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusBadRequest, "unauthrized", nil, nil)
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid order id", nil, nil)
	}

	err = h.svc.CancelOrder(c.Context(), uint(orderID), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "oreder cannot cancel", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Order cancelled successfully", nil, nil)
}
