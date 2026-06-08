package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"gorm.io/gorm"
)

func SetUpRouter(
	  app *fiber.App,
	  db *gorm.DB,
	  userH *handler.UserHandler,
	  AdminH *handler.AdminHandler,
	  cartH *handler.CartHandler,
	  wishH *handler.WishlistHandler,
	  orderH *handler.OrderHandler,
	  addressH *handler.AddressHandler,
	  ) {
	app.Use(logger.New())	

	api := app.Group("/api")
	v1 := api.Group("/v1")
	SetUpUserRouter(v1, db, userH)
	SetUpAdminRouter(v1, db, AdminH)
	SetUpCartRouter(v1, db, cartH)
	SetupWishlistRoutes(v1, db, wishH)
	SetupOrderRoutes(v1, db, orderH)
	SetupAddressRoutes(v1, db, addressH)
}
