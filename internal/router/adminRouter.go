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
		routes.Patch("/production/:id",adminH.UpdateStatus)
		routes.Post("/block/:id",adminH.BlockUser)
		routes.Post("/product",adminH.AddNewProduct)
		routes.Get("/allproducts",adminH.GetAll)
		routes.Get("/Product/:id",adminH.GetProduct)
		routes.Get("/production/:status",adminH.Production)
		routes.Get("/orders", adminH.GetAllOrders)
		routes.Delete("/Product/:id",adminH.DeleteProduct)
		routes.Delete("/user",adminH.DeleteUser)
		routes.Put("/orders/:id", adminH.UpdateStatus)
		
	}
	
}