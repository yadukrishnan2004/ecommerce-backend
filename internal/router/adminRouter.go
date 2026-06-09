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
		admin.Put(AdminUpdateUser, adminH.UpdateUser)
		admin.Patch(AdminUpdateProdStatus, adminH.UpdateStatus)
		admin.Patch(AdminBlockUser, adminH.BlockUser)
		admin.Post(AdminAddProduct, adminH.AddNewProduct)
		admin.Put(AdminUpdateProduct, adminH.UpdateProduct)
		admin.Get(AdminGetAllProducts, adminH.GetAll)
		admin.Get(AdminGetProductsStatus, adminH.Production)
		admin.Get(AdminSearchProducts, adminH.SearchProducts)
		admin.Get(AdminFilterProducts, adminH.FilterProducts)
		admin.Get(AdminGetProduct, adminH.GetProduct)
		admin.Get(AdminGetAllOrders, adminH.GetAllOrders)
		admin.Get(AdminGetOrderDetails, adminH.GetOrderDetails)
		admin.Delete(AdminDeleteProduct, adminH.DeleteProduct)
		admin.Delete(AdminDeleteUser, adminH.DeleteUser)
		admin.Put(AdminUpdateOrderStatus, adminH.UpdateStatus)
		admin.Put(AdminUpdateOrderStatus2, adminH.UpdateOrdersStatus)
		admin.Get(AdminSearchUsers, adminH.SearchUsers)
		admin.Get(AdminGetAllUsers, adminH.GetAllUsers)
		admin.Get(AdminDashboardGraphs, adminH.GetDashboardGraphs)
		admin.Get(AdminGetUser, adminH.GetUser)
		admin.Get(AdminGetUserCart, adminH.GetUserCart)
		admin.Get(AdminGetUserWishlist, adminH.GetUserWishlist)
		admin.Get(AdminGetUserAddresses, adminH.GetUserAddresses)

		// Categories
		admin.Post(AdminCreateCategory, adminH.CreateCategory)
		admin.Get(AdminGetAllCategories, adminH.GetAllCategories)
		admin.Put(AdminUpdateCategory, adminH.UpdateCategory)
		admin.Delete(AdminDeleteCategory, adminH.DeleteCategory)

		// Inventory Alerts
		admin.Get(AdminLowStock, adminH.GetLowStockProducts)

		// KPI Analytics
		admin.Get(AdminDashboardKpis, adminH.GetDashboardKPIs)
	}

}
