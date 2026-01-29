package handler

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/pkg"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/service"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/constants"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/response"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}



func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	var request dto.CreateUserRequest

	if err := c.BodyParser(&request); err != nil {
		return response.Response(c,http.StatusBadRequest,"invalid input",nil,err.Error())
	}

	if err:=pkg.Validate.Struct(request);err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",request, err.Error())
	}

	token, err := h.svc.SignUp(c.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "please try again later",nil, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwtverify",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	return response.Response(c, http.StatusOK, fmt.Sprintf("otp is send to your %s", request.Email),request,nil)
}

func (h *UserHandler) OtpVerify(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return response.Response(c,http.StatusUnauthorized,"unauthorized",email,nil)
	}
	var otp dto.Otp

	if err := c.BodyParser(&otp); err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",otp, err.Error())
	}
	if err:=pkg.Validate.Struct(otp);err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",otp, err.Error())
	}

	if err := h.svc.VerifyOtp(c.Context(), email, otp.Otp); err != nil {
		return response.Response(c, http.StatusNotAcceptable, "invalid code",otp,err.Error())
	}
	return response.Response(c,http.StatusCreated, "user created",nil,nil)
}

func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	var request dto.SignInRequest

	if err := c.BodyParser(&request); err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",request,err.Error())
	}
	if err:=pkg.Validate.Struct(request);err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",request, err.Error())
	}

	token, err := h.svc.SignIn(c.Context(), request.Email, request.Password)
	if err != nil {
		return response.Response(c,http.StatusUnauthorized, "user not found",request,err.Error())
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

	return response.Response(c,http.StatusOK, "login successful",request,nil)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return response.Response(c,http.StatusOK, "logged out successfully",nil,nil)
}

func (h *UserHandler) Forgotpassword(c *fiber.Ctx) error {
	var getemail dto.Getemail
	if err := c.BodyParser(&getemail); err != nil {
		return response.Response(c, constants.BADREQUEST, "email is not valid",nil,err.Error())
	}
	if err:=pkg.Validate.Struct(getemail);err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",getemail,err.Error())
	}

	token, err := h.svc.Forgetpassword(c.Context(), getemail.Email)
	if err != nil {
		return response.Response(c,http.StatusInternalServerError, "please try again later",getemail,err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "forgetpassword",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)
	return response.Response(c,http.StatusOK, fmt.Sprintf("otp is send to your %s", getemail.Email),nil,nil)
}

// =======================================================================================
func (h *UserHandler) Resetpassword(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return response.Response(c,http.StatusUnauthorized, "Unauthorized",nil,nil)
	}
	
	var Reset dto.Reset
	if err := c.BodyParser(&Reset); err != nil {
		return response.Response(c,http.StatusBadGateway,"invalid input",Reset,err.Error())
	}

	if err:=pkg.Validate.Struct(Reset);err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid request",Reset,err.Error())
	}

	if err := h.svc.Resetpassword(c.Context(), email, Reset.Code, Reset.Newpassword); err != nil {
		return response.Response(c,http.StatusInternalServerError, "something went wrong with reset password",nil,err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "forgetpassword",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return response.Response(c,http.StatusOK, "user updated",nil,nil)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("userid").(float64)

	var req dto.UpdateUser
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c,http.StatusBadRequest, "invalid input",nil,err.Error())
	}

	if err := h.svc.UpdateProfile(c.Context(),uint(user), req); err != nil {
		return response.Response(c,http.StatusInternalServerError,"user not updated",req,err.Error())
	}
	return response.Response(c,http.StatusOK, "user updated",nil,nil)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error{
    userIDFloat, ok := c.Locals("userid").(float64)
    if !ok {
        return response.Response(c,http.StatusUnauthorized,"unauthorized",nil,nil)
    }

    user, err := h.svc.GetProfile(c.Context(), uint(userIDFloat))
    if err != nil {
        return response.Response(c,http.StatusUnauthorized,"unauthorized",nil,err.Error())
    }
    return response.Response(c,http.StatusOK,user.Role,user,nil)
}

