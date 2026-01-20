package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type UserHandler struct {
	svc domain.UserService
	vrify domain.UserService
}

func NewUserHandler(svc,vrify domain.UserService) *UserHandler{
	return &UserHandler{
		svc:svc,
		vfy:vrify,
	}
}

func(h *UserHandler) OtpVerify(c *fiber.Ctx)error{
	var otp struct{
		otp string
	}

	if err:=c.BodyParser(&otp);err != nil {
		return c.Status(401).JSON(fiber.Map{"error":"invalid input"})
	}
	
}

func (h *UserHandler) Register(c *fiber.Ctx) error{
	var User struct{
		Name 		string `json:"name"`
		Email		string `json:"email"`
		Password	string `json:"password"`
	}

	if err:=c.BodyParser(&User);err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid input"})
	}

	err:=h.svc.Register(c.Context(),User.Name,User.Email,User.Password)
	if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "An OTP sent to your gmail id"})
}