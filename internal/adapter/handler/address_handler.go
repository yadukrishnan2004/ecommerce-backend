package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

type AddressHandler struct {
	svc usecase.AddressUsecase
}

func NewAddressHandler(svc usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func (h *AddressHandler) AddAddress(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthrized", nil, nil)
	}

	var address domain.Address
	if err := c.BodyParser(&address); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", address, err.Error())
	}

	address.UserID = uint(userIDFloat)

	if err := h.svc.AddAddress(c.Context(), &address); err != nil {
		return response.Response(c, http.StatusInternalServerError, "address not created", nil, err.Error())
	}
	return response.Response(c, http.StatusInternalServerError, "Address added successfully", address, nil)
}

func (h *AddressHandler) GetAddresses(c *fiber.Ctx) error {
	userIDFloat, _ := c.Locals("userid").(float64)

	addresses, err := h.svc.GetUserAddresses(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "address not get", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "data", addresses, nil)
}
