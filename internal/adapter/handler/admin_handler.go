package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/pkg"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type AdminHandler struct {
	svc usecase.AdminUseCase
}

func NewAdminHandler(Svc usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		svc: Svc,
	}
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, err.Error())
	}

	var req dto.AdminUpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	// Map DTO to Input
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

	msg, err := h.svc.BlockUser(c.Context(), uint(id))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, msg, nil, err.Error())
	}
	return response.Response(c, http.StatusOK, msg, nil, nil)
}

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
	product, err := h.svc.GetAllProducts(c.Context())
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "faile to fetch the products", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "all users list", product, nil)
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

func (h *AdminHandler) Production(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return response.Response(c, http.StatusBadRequest, "bad request", nil, "no status found")
	}

	products, err := h.svc.Production(c.Context(), status)
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
	var status dto.UpdateStatus
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

func (h *AdminHandler) GetAllOrders(c *fiber.Ctx) error {

	orders, err := h.svc.GetAllOrders(c.Context())
	if err != nil {
		return response.Response(c, fiber.StatusInternalServerError, "Failed to fetch orders", nil, nil)
	}

	allorders := dto.Orders{
		Count: len(orders),
		Items: orders,
	}

	return response.Response(c, fiber.StatusOK, "get Order successfully", allorders, nil)
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

func (h *AdminHandler) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")

	if query == "" {
		return response.Response(c, http.StatusBadRequest, "Search query is required", nil, nil)
	}

	products, err := h.svc.SearchProducts(c.Context(), query)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Search failed", products, err)
	}

	return response.Response(c, http.StatusOK, "search result", products, nil)
}

func (h *AdminHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.Response(c, http.StatusBadRequest, "Search query is required", nil, nil)
	}

	users, err := h.svc.SearchUsers(c.Context(), query)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Search query is required", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "list of users found", users, nil)
}

func (h *AdminHandler) FilterProducts(c *fiber.Ctx) error {
	filter := domain.ProductFilter{
		Search:   c.Query("search"),
		MinPrice: c.QueryFloat("min_price", 0),
		MaxPrice: c.QueryFloat("max_price", 0),
		Sort:     c.Query("sort"),
		Category: c.Query("category"),
	}

	products, err := h.svc.FilterProducts(c.Context(), filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Failed to fetch products", nil, nil)
	}

	return response.Response(c, http.StatusOK, "success", products, nil)
}
