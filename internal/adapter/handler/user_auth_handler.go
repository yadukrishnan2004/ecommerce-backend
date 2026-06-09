package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/pkg"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	var request dto.CreateUserRequest

	if err := c.BodyParser(&request); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	if err := pkg.Validate.Struct(request); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", request, err.Error())
	}

	if strings.Contains(request.Name, " ") {
		return response.Response(c, http.StatusBadRequest, "name cannot contain spaces", request.Name, "invalid name")
	}

	if strings.ContainsAny(request.Name, ".*#!+=-()/<>,:;'") {
		return response.Response(c, http.StatusBadRequest, "name contains invalid characters .*#!+=-()/<>,:;'", request.Name, "invalid name")
	}

	request.Password = strings.TrimSpace(request.Password)
	if strings.Contains(request.Password, " ") {
		return response.Response(c, http.StatusBadRequest, "password cannot contain spaces", request.Password, "invalid password")
	}

	if strings.ContainsAny(request.Password, ".*#!+=-()/<>,:;'@") {
		return response.Response(c, http.StatusBadRequest, "password contains invalid characters ", ".*#!+=-()/<>,:;'", "invalid password")
	}

	msg, err := h.svc.SignUp(c.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "please try again later", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, msg, nil, nil)
}

func (h *UserHandler) OtpVerify(c *fiber.Ctx) error {
	var otp dto.Otp

	if err := c.BodyParser(&otp); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", otp, err.Error())
	}
	if err := pkg.Validate.Struct(otp); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", otp, err.Error())
	}

	token, err := h.svc.VerifyOtp(c.Context(), otp.Email, otp.Otp)
	if err != nil {
		return response.Response(c, http.StatusNotAcceptable, "invalid code", otp, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	}

	c.Cookie(&cookie)
	return response.Response(c, http.StatusOK, "Login successful", nil, nil)
}

func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	var request dto.SignInRequest

	if err := c.BodyParser(&request); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", request, err.Error())
	}

	if err := pkg.Validate.Struct(request); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", request, err.Error())
	}

	if strings.ContainsAny(request.Password, ".*#!+=-()/<>,:;'") {
		return response.Response(c, http.StatusBadRequest, "password contains invalid characters", nil, "invalid password")
	}

	userdata, rtoken, token, err := h.svc.SignIn(c.Context(), request.Email, request.Password)
	if err != nil {
		return response.Response(c, http.StatusUnauthorized, "user not found", request, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	}

	c.Cookie(&cookie)

	refresh := fiber.Cookie{
		Name:     "refresh",
		Value:    rtoken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	}

	c.Cookie(&refresh)

	return response.Response(c, http.StatusOK, "login successful", userdata, nil)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	}

	c.Cookie(&cookie)
	refresh := fiber.Cookie{
		Name:     "refresh",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	}

	c.Cookie(&refresh)

	return response.Response(c, http.StatusOK, "logged out successfully", nil, nil)
}

func (h *UserHandler) Forgotpassword(c *fiber.Ctx) error {
	var getemail dto.Getemail
	if err := c.BodyParser(&getemail); err != nil {
		return response.Response(c, http.StatusBadRequest, "email is not valid", nil, err.Error())
	}
	if err := pkg.Validate.Struct(getemail); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", getemail, err.Error())
	}

	if !strings.HasSuffix(getemail.Email, "@gmail.com") {
		return response.Response(c, http.StatusBadRequest, "only gmail.com emails allowed", nil, "invalid email domain")
	}

	token, err := h.svc.ForgotPassword(c.Context(), getemail.Email)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "please try again later", getemail, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwtverify",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	}
	c.Cookie(&cookie)
	return response.Response(c, http.StatusOK, fmt.Sprintf("otp is send to your %s", getemail.Email), nil, nil)
}

func (h *UserHandler) Resetpassword(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "Unauthorized", nil, nil)
	}

	var Reset dto.Reset
	if err := c.BodyParser(&Reset); err != nil {
		return response.Response(c, http.StatusBadGateway, "invalid input", Reset, err.Error())
	}

	if err := pkg.Validate.Struct(Reset); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid request", Reset, err.Error())
	}

	if strings.ContainsAny(Reset.Newpassword, ".*#!+=-()/<>,:;'") {
		return response.Response(c, http.StatusBadRequest, "password contains invalid characters", nil, "invalid password")
	}

	if err := h.svc.ResetPassword(c.Context(), email, Reset.Code, Reset.Newpassword); err != nil {
		return response.Response(c, http.StatusInternalServerError, "something went wrong with reset password", nil, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwtverify",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	}

	c.Cookie(&cookie)
	return response.Response(c, http.StatusOK, "user updated", nil, nil)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, nil)
	}

	var req dto.UpdateUser
	if err := c.BodyParser(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, "invalid input", nil, err.Error())
	}

	if req.Name == nil || strings.Contains(*req.Name, " ") {
		return response.Response(c, http.StatusBadRequest, "name cannot contain spaces or be empty", nil, "invalid name")
	}

	if strings.ContainsAny(*req.Name, ".*#!+=-()/<>,:;'") {
		return response.Response(c, http.StatusBadRequest, "name contains invalid characters .*#!+=-()/<>,:;'", req.Name, "invalid name")
	}

	input := usecase.UpdateUserInput{
		Name: req.Name,
	}

	if err := h.svc.UpdateProfile(c.Context(), uint(user), input); err != nil {
		return response.Response(c, http.StatusInternalServerError, "user not updated", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, "user updated", nil, nil)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("userid").(float64)
	if !ok {
		return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, nil)
	}

	user, err := h.svc.GetProfile(c.Context(), uint(userIDFloat))
	if err != nil {
		return response.Response(c, http.StatusUnauthorized, "unauthorized", nil, err.Error())
	}
	return response.Response(c, http.StatusOK, user.Role, user, nil)
}
