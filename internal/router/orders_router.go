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
		order.Post("/", orderH.PlaceOrder)
		order.Get("/", orderH.GetOrderHistory)
		order.Get("/:id", orderH.GetOrder)
		order.Post("/buy-now", orderH.BuyNow)
		order.Post("/verify-payment", orderH.VerifyPayment)
	}
}
