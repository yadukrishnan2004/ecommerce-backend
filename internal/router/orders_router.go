package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetupOrderRoutes(api fiber.Router, db *gorm.DB, orderH *handler.OrderHandler) {
	order := api.Group("/orders")
	order.Use(middleware.UserMiddleware(db))
	{
		order.Post(OrderPlace, orderH.PlaceOrder)
		order.Get(OrderGet, orderH.GetOrderHistory)
		order.Get(OrderGetByID, orderH.GetOrder)
		order.Post(OrderBuyNow, orderH.BuyNow)
		order.Post(OrderVerifyPayment, orderH.VerifyPayment)
	}
}
