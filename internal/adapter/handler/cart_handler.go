package handler

import (
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type CartHandler struct {
	svc usecase.CartService
}

func NewCartHandler(svc usecase.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {

	var req dto.CartRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input data", req, err.Error())
	}
	if req.Quantity <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid input", req, "Quantity must be greater than 0")
	}
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "no user found", nil, nil)
	}
	userID := uint(userIDFloat)
	err := h.svc.AddToCart(c.Context(), userID, req.ProductID, uint(req.Quantity))
	if err != nil {
		if err.Error() == "product not found" {
			return response.Response(c, http.StatusNotFound, "product not found", req, err.Error())
		}
		if err.Error() == "insufficient stock" {
			return response.Response(c, http.StatusBadRequest, "insufficient stock", req, err.Error())
		}
		return response.Response(c, http.StatusInternalServerError, "cart service not working", req, err.Error())
	}

	return response.Response(c, http.StatusOK, "item added to cart", req, nil)
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "no user found", nil, nil)
	}
	err := h.svc.ClearCart(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "faile to clear the cart", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Cart cleared successfully", nil, nil)
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "no user found", nil, nil)
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid product id", nil, err.Error())
	}

	err = h.svc.RemoveItem(c.Context(), uint(userIDFloat), uint(productID))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "item not removed from cart", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "item remover successfully", nil, nil)
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "no user found", nil, nil)
	}

	items, err := h.svc.GetCart(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "not get any cart item", nil, err.Error())
	}

	var grandTotal float64

	var responseItems []dto.CartItemResponse

	for _, item := range items {
		subTotal := float64(item.Product.Price) * float64(item.Quantity)
		grandTotal += subTotal

		responseItems = append(responseItems, dto.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       float64(item.Product.Price),
			Quantity:    int(item.Quantity),
			SubTotal:    subTotal,
		})
	}

	cart := dto.Cart{
		Items:      responseItems,
		GrandTotal: float32(grandTotal),
		Count:      uint(len(items)),
	}

	return response.Response(c, http.StatusOK, "get cart successfully", cart, nil)

}

func (h *CartHandler) UpdateQuantity(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "no user found", nil, nil)
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid product id", nil, err.Error())
	}

	var req dto.UpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input data", req, err.Error())
	}

	err = h.svc.UpdateQuantity(c.Context(), uint(userIDFloat), uint(productID), req.Quantity)
	if err != nil {
		if err.Error() == "product not found" {
			return response.Response(c, http.StatusNotFound, "product not found", nil, err.Error())
		}
		if err.Error() == "insufficient stock available" {
			return response.Response(c, http.StatusBadRequest, "insufficient stock", nil, err.Error())
		}
		if err.Error() == "quantity must be greater than 0" {
			return response.Response(c, http.StatusBadRequest, "invalid quantity", nil, err.Error())
		}
		return response.Response(c, http.StatusInternalServerError, "item not updated", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "quantity updated", nil, nil)
}
