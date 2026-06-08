package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/pkg"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *AdminHandler) AddNewProduct(c *fiber.Ctx) error {
	var newProduct domain.Product

	if err := c.BodyParser(&newProduct); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", newProduct, err.Error())
	}

	if err := pkg.Validate.Struct(newProduct); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", newProduct, err.Error())
	}
	product := domain.Product{
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Description: newProduct.Description,
		Category:    newProduct.Category,
		Offer:       newProduct.Offer,
		OfferPrice:  newProduct.OfferPrice,
		Production:  newProduct.Production,
		Images:      newProduct.Images,
		Stock:       newProduct.Stock,
	}
	if err := h.svc.AddNewProduct(c.Context(), &product); err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed add new product", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "product created", nil, nil)
}

func (h *AdminHandler) GetAll(c *fiber.Ctx) error {
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

func (h *AdminHandler) GetProduct(c *fiber.Ctx) error {
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

func (h *AdminHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	var req domain.Product
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	if err := h.svc.UpdateProduct(c.Context(), uint(id), &req); err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to update product", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "product updated successfully", nil, nil)
}

func (h *AdminHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "no requst found", nil, err.Error())
	}

	if err := h.svc.DeleteProduct(c.Context(), uint(id)); err != nil {
		return response.Response(c, http.StatusBadRequest, "Product not found", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "product deleted Successfully", nil, nil)
}

func (h *AdminHandler) Production(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return response.Response(c, http.StatusBadRequest, "bad request", nil, "no status found")
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

	products, err := h.svc.Production(c.Context(), status, limit, offset)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "not found", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "product get successfully", products, nil)
}

func (h *AdminHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, "bad request", nil, "no status found")
	}
	var status struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&status); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}
	if err := pkg.Validate.Struct(status); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}
	if err := h.svc.UpdateStatus(c.Context(), uint(id), status.Status); err != nil {
		return response.Response(c, http.StatusInternalServerError, "try again later", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "status updated", nil, nil)
}

func (h *AdminHandler) SearchProducts(c *fiber.Ctx) error {
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

func (h *AdminHandler) FilterProducts(c *fiber.Ctx) error {
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

func (h *AdminHandler) GetLowStockProducts(c *fiber.Ctx) error {
	threshold, _ := strconv.Atoi(c.Query("threshold", "5"))

	products, err := h.svc.GetLowStockProducts(c.Context(), threshold)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to fetch low stock products", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "low stock products fetched successfully", products, nil)
}
