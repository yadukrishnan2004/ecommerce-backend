package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetupAddressRoutes(api fiber.Router, addressH *handler.AddressHandler) {
	addr := api.Group("/addresses")
	addr.Use(middleware.UserMiddleware)

	addr.Post("/", addressH.AddAddress) 
	addr.Get("/", addressH.GetAddresses) 
}