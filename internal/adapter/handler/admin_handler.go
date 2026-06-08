package handler

import (
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
)

type AdminHandler struct {
	svc usecase.AdminUseCase
}

func NewAdminHandler(Svc usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		svc: Svc,
	}
}
