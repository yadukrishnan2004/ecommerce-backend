package handler

import (
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
)

type UserHandler struct {
	svc usecase.UserUseCase
}

func NewUserHandler(svc usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}
