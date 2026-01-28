package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/constants"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utile/response"
)


type AdminHandler struct{
	svc domain.AdminService
}

func NewAdminHandler(Svc domain.AdminService) *AdminHandler{
	return &AdminHandler{
		svc:Svc,
	}
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx)error{
	var req struct{
	Id		  uint	 `json:"id" binding:"required"`
	Name 	  string `json:"name"`
	Email     string `json:"email"`
	Role 	  string `json:"role"`
	Otp		  string `json:"otp"`
	IsActive  bool   `json:"is_active"`
	IsBlocked bool   `json:"is_blocked"` 
	OtpExpire int64  `json:"otp_expire"`
	}

	if err:=c.BodyParser(&req);err !=nil{
		return response.Error(c,constants.BADREQUEST,"invalid input",err)
	}
	if req.Id <=0{
		return response.Error(c,constants.BADREQUEST,"invalid input","no id provided")
	}
	if req.Name==""{
		return response.Error(c,constants.BADREQUEST,"invalid input","invalid input at name")
	}
	if req.Email == ""{
		return response.Error(c,constants.BADREQUEST,"invalid input","invalid input at email")
	}


		user:=&domain.User{
		Name: req.Name,
		Email: req.Email,
		Role: req.Role,
		Otp: req.Otp,
		IsActive: req.IsActive,
		IsBlocked: req.IsBlocked,
		OtpExpire: req.OtpExpire,
	}
	updateuser,err:=h.svc.UpdateUser(c.Context(),req.Id,*user)
	if err != nil {
		return response.Error(c,constants.INTERNALSERVERERROR,"user not updated",err.Error())
	}
	return response.Success(c,constants.SUCCESS,"user updated",updateuser)
}

func (h *AdminHandler) BlockUser(c *fiber.Ctx) error{
	var user struct{
		Id  uint `json:"id" binding:"required"`
	}
	if err:=c.BodyParser(&user); err != nil{
		return response.Error(c,constants.BADREQUEST,"user id required for block",err.Error())
	}

	if err:=h.svc.BlockUser(c.Context(),user.Id);err != nil{
		return response.Error(c,constants.BADREQUEST,"service error",err.Error())
	}

	return response.Success(c,constants.SUCCESS,"user blocked successfully","")
}

