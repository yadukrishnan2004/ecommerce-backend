package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

type UserHandler struct {
	svc domain.UserService

}

func NewUserHandler(svc domain.UserService) *UserHandler{
	return &UserHandler{
		svc:svc,
	}
}

func(h *UserHandler) OtpVerify(c *fiber.Ctx)error{
	var otp struct{
		Email string `json:"email"`
		Otp string	`json:"otp"`
	}

	if err:=c.BodyParser(&otp);err != nil {
		return c.Status(401).JSON(fiber.Map{"error":"invalid input"})
	}

	if err:=h.svc.VerifyOtp(c.Context(),otp.Email,otp.Otp);err != nil {
		return c.Status(401).JSON(fiber.Map{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Account verified successfully","status":"user created"})

	
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

func (h *UserHandler) Login(c *fiber.Ctx) error{
	type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

	token, err := h.svc.Login(c.Context(), req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }
	//  Create the Cookie (HTTP Logic)
    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(15 * time.Minute), // Match your token expiry
        HTTPOnly: true,                           // XSS Protection: JS cannot read this
        Secure:   false,                          // Set to TRUE in production (HTTPS)
        SameSite: "Lax",                          // CSRF Protection
    }

    //  Attach cookie to response
    c.Cookie(&cookie)


    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
    })
}