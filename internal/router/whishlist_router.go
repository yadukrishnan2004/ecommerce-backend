package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetupWishlistRoutes(api fiber.Router, wishH *handler.WishlistHandler) {
	wish := api.Group("/wishlist")
	wish.Use(middleware.UserMiddleware)
	{
		wish.Post("/:id", wishH.AddToWishlist)        
		wish.Delete("/:id", wishH.RemoveFromWishlist) 
		wish.Delete("/clear", wishH.ClearWishlist)
		wish.Get("/", wishH.GetWishlist)
	}
}
