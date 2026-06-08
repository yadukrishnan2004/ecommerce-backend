package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	var req dto.AdminUpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	input := usecase.AdminUpdateUserInput{
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		IsActive:  req.IsActive,
		IsBlocked: req.IsBlocked,
	}

	updateuser, err := h.svc.UpdateUser(c.Context(), uint(id), input)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "user not updated", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "user updated", updateuser, nil)
}

func (h *AdminHandler) BlockUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	var body struct {
		Blocked *bool `json:"blocked"`
	}

	if err := c.BodyParser(&body); err != nil && err != fiber.ErrUnprocessableEntity {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	var blockedOpt *bool
	if body.Blocked != nil {
		blockedOpt = body.Blocked
	}

	msg, err := h.svc.BlockUser(c.Context(), uint(id), blockedOpt)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, msg, nil, err.Error())
	}
	return response.Response(c, http.StatusOK, msg, nil, nil)
}

func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "no requst found", nil, err.Error())
	}

	if err := h.svc.DeleteUser(c.Context(), uint(id)); err != nil {
		return response.Response(c, http.StatusBadRequest, "user not found", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "user Deleted Successfully", nil, nil)
}

func (h *AdminHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.Response(c, http.StatusBadRequest, "Search query is required", nil, nil)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, err := h.svc.SearchUsers(c.Context(), query, limit, offset)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Search failed", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "list of users found", users, nil)
}

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, err := h.svc.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return response.Response(c, http.StatusOK, "Get all users Success", users, nil)
}

func (h *AdminHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}
	user, err := h.svc.GetUserByID(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusNotFound, "user not found", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "user fetched successfully", user, nil)
}

func (h *AdminHandler) GetUserCart(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}
	cart, err := h.svc.GetUserCart(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to get cart", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "cart fetched successfully", cart, nil)
}

func (h *AdminHandler) GetUserWishlist(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}
	wishlist, err := h.svc.GetUserWishlist(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to get wishlist", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "wishlist fetched successfully", wishlist, nil)
}

func (h *AdminHandler) GetUserAddresses(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}
	addresses, err := h.svc.GetUserAddresses(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to get addresses", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "addresses fetched successfully", addresses, nil)
}
