package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *UserHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "no requst found", nil, err.Error())
	}
	product, err := h.svc.GetProduct(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "Product not found", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "product found", product, nil)
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	product, err := h.svc.GetAllProducts(c.Context(), limit, offset)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "faile to fetch the products", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "all products list", product, nil)
}

func (h *UserHandler) SearchProducts(c *fiber.Ctx) error {
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

	products, err := h.svc.SearchProducts(c.Context(), query, limit, offset)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Search failed", products, err)
	}

	return response.Response(c, http.StatusOK, "search result", products, nil)
}

func (h *UserHandler) FilterProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	filter := domain.ProductFilter{
		Search:   c.Query("search"),
		MinPrice: c.QueryFloat("min_price", 0),
		MaxPrice: c.QueryFloat("max_price", 0),
		Sort:     c.Query("sort"),
		Category: c.Query("category"),
		Limit:    limit,
		Offset:   offset,
	}

	products, err := h.svc.FilterProducts(c.Context(), filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Failed to fetch products", nil, nil)
	}

	return response.Response(c, http.StatusOK, "success", products, nil)
}
