package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetUpCartRouter(api fiber.Router, db *gorm.DB, cartH *handler.CartHandler){
	cart:=api.Group("/cart")
	routes:=cart.Group("/")
	routes.Use(middleware.UserMiddleware(db))
	{
		routes.Post(CartAdd, cartH.AddToCart)
		routes.Delete(CartClear, cartH.ClearCart)
		routes.Delete(CartRemove, cartH.RemoveItem)
		routes.Get(CartGet, cartH.GetCart)
		routes.Put(CartUpdate, cartH.UpdateQuantity)
	}
}