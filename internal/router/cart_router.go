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
		routes.Post("/add",cartH.AddToCart)
		routes.Delete("/clear", cartH.ClearCart)
		routes.Delete("/:id", cartH.RemoveItem)
		routes.Get("/",cartH.GetCart)
		routes.Put("/:id", cartH.UpdateQuantity)
	}
}