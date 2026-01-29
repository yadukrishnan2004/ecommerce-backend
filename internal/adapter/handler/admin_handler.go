package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
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
