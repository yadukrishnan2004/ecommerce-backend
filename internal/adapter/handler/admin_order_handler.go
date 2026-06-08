package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/pkg"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *AdminHandler) GetAllOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	orders, err := h.svc.GetAllOrders(c.Context(), limit, offset)
	if err != nil {
		return response.Response(c, fiber.StatusInternalServerError, "Failed to fetch orders", nil, nil)
	}

	allorders := dto.Orders{
		Count: len(orders),
		Items: orders,
	}

	return response.Response(c, fiber.StatusOK, "get Order successfully", allorders, nil)
}

func (h *AdminHandler) GetOrderDetails(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil || orderID <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid order id", nil, nil)
	}

	orderItems, err := h.svc.GetOrderDetails(c.Context(), uint(orderID))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to get order details", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Order details fetched successfully", orderItems, nil)
}

func (h *AdminHandler) UpdateOrdersStatus(c *fiber.Ctx) error {
	orderID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid orderId", nil, nil)
	}
	var req dto.UpdateStatus
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, nil)
	}

	if err := pkg.Validate.Struct(req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	err = h.svc.UpdateOrderStatus(c.Context(), uint(orderID), req.Status)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "status not updated", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "order status updated", nil, nil)
}
