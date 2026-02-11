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
		admin.Patch("/users/:id", adminH.UpdateUser)             // Changed from /update/:id to /users/:id (RESTful)
		admin.Patch("/products/:id/status", adminH.UpdateStatus) // Changed from /production/:id
		admin.Post("/users/:id/block", adminH.BlockUser)         // Changed from /block/:id
		admin.Post("/products", adminH.AddNewProduct)            // Changed from /product
		admin.Get("/products", adminH.GetAll)                    // Changed from /allproducts
		admin.Get("/products/:id", adminH.GetProduct)            // Fixed casing
		admin.Get("/products/status/:status", adminH.Production) // Changed from /production/:status
		admin.Get("/orders", adminH.GetAllOrders)
		admin.Delete("/products/:id", adminH.DeleteProduct) // Fixed casing
		admin.Delete("/users/:id", adminH.DeleteUser)       // Changed from /user
		admin.Put("/orders/:id", adminH.UpdateStatus)
		admin.Get("/products/search", adminH.SearchProducts) // Changed from /search
		admin.Get("/users/search", adminH.SearchUsers)
		admin.Get("/products/filter", adminH.FilterProducts) // Changed from /filter
	}

}
