package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetupAddressRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config, addressH *handler.AddressHandler) {
	addr := api.Group("/addresses")
	addr.Use(middleware.UserMiddleware(db, cfg))

	addr.Post(AddressPost, addressH.AddAddress)
	addr.Get(AddressGet, addressH.GetAddresses)
}