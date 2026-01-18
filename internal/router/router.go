package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
)

func SetUpRouther(app *fiber.App,userH *handler.UserHandler){
	app.Use(logger.New())

	api:=app.Group("/api")
	v1:=api.Group("/v1")

	UserRouter:=v1.Group("/user")
	{
		UserRouter.Post("/signup",userH.Register)
	}
	
}