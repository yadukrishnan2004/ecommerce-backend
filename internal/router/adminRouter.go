package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetUpAdminRouter(api fiber.Router, db *gorm.DB, adminH *handler.AdminHandler) {
	admin := api.Group("/admin")
	admin.Use(middleware.Adminmiddleware(db))
	{
		admin.Put("/users/:id", adminH.UpdateUser)
		admin.Patch("/products/:id/status", adminH.UpdateStatus)
		admin.Patch("/users/:id/block", adminH.BlockUser)
		admin.Post("/products", adminH.AddNewProduct)
		admin.Put("/products/:id", adminH.UpdateProduct)
		admin.Get("/products", adminH.GetAll)
		admin.Get("/products/status/:status", adminH.Production)
		admin.Get("/products/search", adminH.SearchProducts)
		admin.Get("/products/filter", adminH.FilterProducts)
		admin.Get("/products/:id", adminH.GetProduct)
		admin.Get("/orders", adminH.GetAllOrders)
		admin.Get("/orders/:id", adminH.GetOrderDetails)
		admin.Delete("/products/:id", adminH.DeleteProduct)
		admin.Delete("/users/:id", adminH.DeleteUser)
		admin.Put("/orders/:id", adminH.UpdateStatus)
		admin.Put("/orders/status/:id", adminH.UpdateOrdersStatus)
		admin.Get("/users/search", adminH.SearchUsers)
		admin.Get("/users", adminH.GetAllUsers)
		admin.Get("/dashboard-graphs", adminH.GetDashboardGraphs)
		admin.Get("/users/:id", adminH.GetUser)
		admin.Get("/users/:id/cart", adminH.GetUserCart)
		admin.Get("/users/:id/wishlist", adminH.GetUserWishlist)
		admin.Get("/users/:id/addresses", adminH.GetUserAddresses)

		// Categories
		admin.Post("/categories", adminH.CreateCategory)
		admin.Get("/categories", adminH.GetAllCategories)
		admin.Put("/categories/:id", adminH.UpdateCategory)
		admin.Delete("/categories/:id", adminH.DeleteCategory)

		// Inventory Alerts
		admin.Get("/inventory/low-stock", adminH.GetLowStockProducts)

		// KPI Analytics
		admin.Get("/dashboard/kpis", adminH.GetDashboardKPIs)
	}

}
