package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetupOrderRoutes(api fiber.Router, orderH *handler.OrderHandler) {
	order := api.Group("/orders")
	order.Use(middleware.UserMiddleware)
	{
		order.Post("/", orderH.PlaceOrder)
		order.Get("/", orderH.GetOrderHistory)
		order.Post("/buy-now", orderH.BuyNow)
	}
}
