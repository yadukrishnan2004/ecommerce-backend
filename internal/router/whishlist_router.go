package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetupWishlistRoutes(api fiber.Router, wishH *handler.WishlistHandler) {
	wish:= api.Group("/wishlist")
	routes:= wish.Group("/")
	routes.Use(middleware.UserMiddleware)
	{
		routes.Post("/add/:id",wishH.AddToWishlist)
		routes.Delete("/remove/:id",wishH.RemoveFromWishlist)
		routes.Delete("/", wishH.ClearWishlist)
		routes.Get("/wishlist",wishH.GetWishlist)
	}

}