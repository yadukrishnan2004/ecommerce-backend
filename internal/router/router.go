package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
)

func SetUpRouter(
	  app *fiber.App,
	  userH *handler.UserHandler,
	  AdminH *handler.AdminHandler,
	  cartH *handler.CartHandler,
	  wishH *handler.WishlistHandler,
	  ) {
	app.Use(logger.New())	

	api := app.Group("/api")
	v1 := api.Group("/v1")
	SetUpUserRouter(v1, userH)
	SetUpAdminRouter(v1, AdminH)
	SetUpCartRouter(v1,cartH)
	SetupWishlistRoutes(v1,wishH)
}
