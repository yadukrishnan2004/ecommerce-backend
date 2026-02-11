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
		wish.Post("/:id", wishH.AddToWishlist)        // Changed from /add/:id - POST /wishlist/:id implies adding
		wish.Delete("/:id", wishH.RemoveFromWishlist) // Changed from /remove/:id - DELETE /wishlist/:id implies removing
		wish.Delete("/", wishH.ClearWishlist)
		wish.Get("/", wishH.GetWishlist)
	}
}
