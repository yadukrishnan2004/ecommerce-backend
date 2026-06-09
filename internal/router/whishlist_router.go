package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetupWishlistRoutes(api fiber.Router, db *gorm.DB, wishH *handler.WishlistHandler) {
	wish := api.Group("/wishlist")
	wish.Use(middleware.UserMiddleware(db))
	{
		wish.Post(WishlistAdd, wishH.AddToWishlist)
		wish.Delete(WishlistClear, wishH.ClearWishlist)
		wish.Delete(WishlistRemove, wishH.RemoveFromWishlist)
		wish.Get(WishlistGet, wishH.GetWishlist)
	}
}
