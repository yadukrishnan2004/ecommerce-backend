package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetUpAdminRouter(api fiber.Router,adminH *handler.AdminHandler){
	admin:=api.Group("/admin")
	routes:=admin.Group("/")
	routes.Use(middleware.Adminmiddleware)
	{
		routes.Patch("/update",adminH.UpdateUser)
		routes.Post("/block",adminH.BlockUser)
	}
	
}