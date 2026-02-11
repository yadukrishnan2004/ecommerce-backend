package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type WishlistHandler struct {
	svc usecase.WishlistService
}

func NewWishlistHandler(svc usecase.WishlistService) *WishlistHandler {
	return &WishlistHandler{svc: svc}
}

func (h *WishlistHandler) AddToWishlist(c *fiber.Ctx) error {

	id, erro := strconv.Atoi(c.Params("id"))

	if erro != nil {
		return response.Response(c, http.StatusBadRequest, "Invalid JSON", id, erro.Error())
	}

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusBadRequest, "unauthrized", nil, nil)

	}

	err := h.svc.AddToWishlist(c.Context(), uint(userIDFloat), uint(id))
	if err != nil {

		if err.Error() == "item already in wishlist" {
			return response.Response(c, http.StatusBadRequest, "it already in the whishlist", nil, err.Error())
		}
		return response.Response(c, http.StatusInternalServerError, "something wrong", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "item added to the wishlist", nil, nil)
}

func (h *WishlistHandler) RemoveFromWishlist(c *fiber.Ctx) error {
	id, erro := strconv.Atoi(c.Params("id"))

	if erro != nil {
		return response.Response(c, http.StatusBadRequest, "Invalid JSON", id, erro.Error())
	}

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusBadRequest, "unauthrized", nil, nil)

	}

	err := h.svc.RemoveFromWishlist(c.Context(), uint(userIDFloat), uint(id))
	if err != nil {

		if err.Error() == "wish list don't have an item" {
			return response.Response(c, http.StatusBadRequest, "it is not in  the whishlist", nil, err.Error())
		}
		return response.Response(c, http.StatusInternalServerError, "something wrong", nil, err.Error())
	}

	return response.Response(c, fiber.StatusOK, "item removed from the wishlist", nil, nil)
}

func (h *WishlistHandler) ClearWishlist(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusBadRequest, "unauthrized", nil, nil)
	}

	err := h.svc.ClearWishlist(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Failed to clear wishlist", nil, err.Error())
	}

	return response.Response(c, fiber.StatusOK, "Wishlist cleared successfully", nil, nil)
}

func (h *WishlistHandler) GetWishlist(c *fiber.Ctx) error {

	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, fiber.StatusBadRequest, "unauthrized", nil, nil)
	}

	items, err := h.svc.GetWishlist(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "faile to fetch wishlist", nil, err.Error())
	}
	result := dto.Wishlist{
		Item:  items,
		Count: len(items),
	}

	return response.Response(c, fiber.StatusOK, "get wishlist successfully", result, nil)
}
