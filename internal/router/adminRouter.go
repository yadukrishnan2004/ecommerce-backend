package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetUpAdminRouter(api fiber.Router, adminH *handler.AdminHandler) {
	admin := api.Group("/admin")
	admin.Use(middleware.Adminmiddleware)
	{
		admin.Patch("/users/:id", adminH.UpdateUser)             
		admin.Patch("/products/:id/status", adminH.UpdateStatus) 
		admin.Post("/users/:id/block", adminH.BlockUser)         
		admin.Post("/products", adminH.AddNewProduct)            
		admin.Get("/products", adminH.GetAll)                    
		admin.Get("/products/status/:status", adminH.Production) 
		admin.Get("/products/search", adminH.SearchProducts)     
		admin.Get("/products/filter", adminH.FilterProducts)     
		admin.Get("/products/:id", adminH.GetProduct)            
		admin.Get("/orders", adminH.GetAllOrders)
		admin.Delete("/products/:id", adminH.DeleteProduct) 
		admin.Delete("/users/:id", adminH.DeleteUser)       
		admin.Put("/orders/:id", adminH.UpdateStatus)
		admin.Get("/users/search", adminH.SearchUsers)
	}

}
