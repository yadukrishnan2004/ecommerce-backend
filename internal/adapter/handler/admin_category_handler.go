package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *AdminHandler) CreateCategory(c *fiber.Ctx) error {
	var category domain.Category
	if err := c.BodyParser(&category); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	if err := h.svc.CreateCategory(c.Context(), &category); err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to create category", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "category created successfully", category, nil)
}

func (h *AdminHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.svc.GetAllCategories(c.Context())
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to fetch categories", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "categories fetched successfully", categories, nil)
}

func (h *AdminHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	if err := h.svc.UpdateCategory(c.Context(), uint(id), req.Name, req.Description); err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to update category", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "category updated successfully", nil, nil)
}

func (h *AdminHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	if err := h.svc.DeleteCategory(c.Context(), uint(id)); err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to delete category", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "category deleted successfully", nil, nil)
}
