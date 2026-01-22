package handler

import (
	"time"
    "fmt"
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
	//   Cookie (HTTP Logic)
    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(15 * time.Minute), 
        HTTPOnly: true,                           
        Secure:   false,                          
        SameSite: "Lax",                          
    }

    //  Attach cookie to response
    c.Cookie(&cookie)


    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
    })
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour), 
        HTTPOnly: true,
    }

    c.Cookie(&cookie)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Logged out successfully",
    })
}

func (h *UserHandler) Forgetpassword(c *fiber.Ctx) error{
    var getemail struct{
        Email string `json:"email" binding:"required,email"`
    }
   if err:= c.BodyParser(&getemail);err != nil {
    return c.Status(400).JSON(fiber.Map{"error":"email is not valid"})
   }

    token,err := h.svc.Forgetpassword(c.Context(), getemail.Email);
     if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "error in forget password",
            "detail": err.Error(),
        })
    }

        cookie := fiber.Cookie{
        Name:     "forgotpassword",
        Value:    token,
        Expires:  time.Now().Add(10 * time.Minute), 
        HTTPOnly: true,                           
        Secure:   false,                          
        SameSite: "Lax",                          
    }
    c.Cookie(&cookie)
    return c.Status(200).JSON(fiber.Map{
        "status":fmt.Sprintf("otp is send to your %s",getemail.Email),
    })

}

func(h *UserHandler) Resetpassword(c *fiber.Ctx)error{
   email :=c.Get("email")

    var Reset struct{
        Code        string `json:"code" binding:"required"`
        Newpassword string `json:"password" binding:"required"`
    }

   if err:=c.BodyParser(&Reset);err !=nil{
    return c.Status(400).JSON(fiber.Map{"error":"invalid input"})
   }

  if err:=h.svc.Resetpassword(c.Context(),email,Reset.Code,Reset.Newpassword);err != nil{
    return  c.Status(400).JSON(fiber.Map{"error":err.Error()})
  }
  return c.Status(201).JSON(fiber.Map{"status":"user updated"})
}

