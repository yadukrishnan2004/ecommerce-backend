package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetupAddressRoutes(api fiber.Router, db *gorm.DB, addressH *handler.AddressHandler) {
	addr := api.Group("/addresses")
	addr.Use(middleware.UserMiddleware(db))

	addr.Post("/", addressH.AddAddress) 
	addr.Get("/", addressH.GetAddresses) 
}