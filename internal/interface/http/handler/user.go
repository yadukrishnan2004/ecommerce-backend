package handler

import "github.com/gofiber/fiber/v2"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserCase interface {
	Execute(email, password string) error
}

type UserHandler struct {
	registerUC RegisterUserCase
}

func NewUserHandler(uc RegisterUserCase) *UserHandler {
	return &UserHandler{
		registerUC: uc,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error{
	var req RegisterRequest
	if err:=c.BodyParser(&req);err != nil {
		return  fiber.ErrBadRequest
	}

	if err :=h.registerUC.Execute(req.Email,req.Password);err != nil {
		return fiber.NewError(fiber.StatusBadRequest,err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}