package handler

import "github.com/gofiber/fiber/v2"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserCase interface {
	Execute(email, password string) error
}

// ---------------------------------------------

type LoginRequest struct{
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUseCase interface{
	Execute(email,password string) error
}


// ---------------------------------------------------------

type UserHandler struct {
	registerUC RegisterUserCase
	loginUC LoginUseCase
}

func NewUserHandler(uc RegisterUserCase,loginUc LoginUseCase) *UserHandler {
	return &UserHandler{
		registerUC: uc,
		loginUC: loginUc,
	}
}


// -----------------------------------------------------------------------
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


func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.loginUC.Execute(req.Email, req.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
