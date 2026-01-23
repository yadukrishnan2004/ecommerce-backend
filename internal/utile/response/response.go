package response

import "github.com/gofiber/fiber/v2"

func Success(ctx *fiber.Ctx,status int,message string,data interface{})error{
	resp:=ApiResponse{
		StatusCode: status,
		Message: message,
	}
	if data != nil{
		resp.Data=data
	}
	return ctx.Status(status).JSON(resp)
}

func Error(ctx *fiber.Ctx,status int,message string,err interface{})error{
	resp:=ApiResponse{
		StatusCode: status,
		Message: message,
	}

	if err != nil{
		resp.Error=err
	}
	return ctx.Status(status).JSON(resp)
}