package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/constants"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/response"
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
    email, ok := c.Locals("email").(string)
    if !ok {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }
	var otp struct{
		Otp string	`json:"otp"`
	}

	if err:=c.BodyParser(&otp);err != nil {
        return response.Error(c,constants.BADREQUEST,"invalid request body",err)
	}

	if err:=h.svc.VerifyOtp(c.Context(),email,otp.Otp);err != nil {
        return response.Error(c,constants.BADREQUEST,"invalid code",err)
	}

    return response.Success(c,constants.CREATED,"user created","")
	
}

func (h *UserHandler) Register(c *fiber.Ctx) error{
	var User struct{
		Name 		string `json:"name"`
		Email		string `json:"email"`
		Password	string `json:"password"`
	}

	if err:=c.BodyParser(&User);err != nil {
        return response.Error(c,constants.BADREQUEST,"invalid input",err)
	}

	token,err:=h.svc.Register(c.Context(),User.Name,User.Email,User.Password)
	if err != nil {
        return response.Error(c,constants.INTERNALSERVERERROR,"please try again later",err)
    }

        cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(10 * time.Minute), 
        HTTPOnly: true,                           
        Secure:   false,                          
        SameSite: "Lax",                          
    }
    c.Cookie(&cookie)

    return response.Success(c,constants.SUCCESSSUCCESS,fmt.Sprintf("otp is send to your %s",User.Email),"")
}

func (h *UserHandler) Login(c *fiber.Ctx) error{
	type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
        return response.Error(c,constants.BADREQUEST,"invalid request",err)
    }

	token, err := h.svc.Login(c.Context(), req.Email, req.Password)
    if err != nil {
        return response.Error(c,constants.UNAUTHORIZED,"user not found",err)
    }

    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(15 * time.Minute), 
        HTTPOnly: true,                           
        Secure:   false,                          
        SameSite: "Lax",                          
    }

    c.Cookie(&cookie)

    return response.Success(c,constants.SUCCESSSUCCESS,"login successful","")
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour), 
        HTTPOnly: true,
    }

    c.Cookie(&cookie)

    return response.Success(c,constants.SUCCESSSUCCESS,"logged out successfully","")
}

func (h *UserHandler) Forgetpassword(c *fiber.Ctx) error{
    var getemail struct{
        Email string `json:"email" binding:"required"`
    }
   if err:= c.BodyParser(&getemail);err != nil {
    return response.Error(c,constants.BADREQUEST,"email is not valid",err)
   }

    token,err := h.svc.Forgetpassword(c.Context(), getemail.Email);
     if err != nil {
        return response.Error(c,constants.INTERNALSERVERERROR,"please try again later",err)
    }

        cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(10 * time.Minute), 
        HTTPOnly: true,                           
        Secure:   false,                          
        SameSite: "Lax",                          
    }
    c.Cookie(&cookie)
    return response.Success(c,constants.SUCCESSSUCCESS,fmt.Sprintf("otp is send to your %s",getemail.Email),"")
}

func(h *UserHandler) Resetpassword(c *fiber.Ctx)error{
  email, ok := c.Locals("email").(string)
    if !ok {
        return response.Error(c,constants.UNAUTHORIZED,"Unauthorized","")
    }

    var Reset struct{
        Code        string `json:"code" binding:"required"`
        Newpassword string `json:"password" binding:"required"`
    }

   if err:=c.BodyParser(&Reset);err !=nil{
    return response.Error(c,constants.BADREQUEST,"invalid input","")
   }

  if err:=h.svc.Resetpassword(c.Context(),email,Reset.Code,Reset.Newpassword);err != nil{
    return response.Error(c,constants.BADREQUEST,"something went wrong with reset password",err.Error())
  }
    return response.Error(c,constants.SUCCESSSUCCESS,"user updated","")
}

