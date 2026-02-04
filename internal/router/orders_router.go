package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetupOrderRoutes(api fiber.Router, orderH *handler.OrderHandler) {
	order := api.Group("/orders")
	routes:=order.Group("/")
	routes.Use(middleware.UserMiddleware)
	{
	routes.Post("/", orderH.PlaceOrder)
	routes.Get("/", orderH.GetOrderHistory)
	routes.Post("/buy-now", orderH.BuyNow)
	}
}