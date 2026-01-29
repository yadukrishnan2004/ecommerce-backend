package handler

import (
	"net/http"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler/dto"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/service"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/response"
)


type AdminHandler struct{
	svc service.AdminService
}

func NewAdminHandler(Svc service.AdminService) *AdminHandler{
	return &AdminHandler{
		svc:Svc,
	}
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx)error{

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
    return response.Response(c,http.StatusBadRequest,"invalid id",nil,err.Error())
	}

	var req dto.AdminUpdateUserRequest

	if err:=c.BodyParser(&req);err !=nil{
		return response.Response(c, http.StatusBadRequest, "invalid id", nil, "id must be positive")
	}

	updateuser,err:=h.svc.UpdateUser(c.Context(),uint(id),&req)
	if err != nil {
		return response.Response(c,http.StatusInternalServerError,"user not updated",updateuser,err.Error())
	}
	return response.Response(c,http.StatusOK,"user updated",updateuser,nil)
}

func (h *AdminHandler) BlockUser(c *fiber.Ctx) error{

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
    return response.Response(c,http.StatusBadRequest,"invalid id",nil,err.Error())
	}

	 msg,err:=h.svc.BlockUser(c.Context(),uint(id));
	 if err != nil{
		return response.Response(c,http.StatusInternalServerError,msg,nil,err.Error())
	}
	return response.Response(c,http.StatusOK,msg,nil,nil)
}
