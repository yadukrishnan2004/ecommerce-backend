package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetUpCartRouter(api fiber.Router,cartH *handler.CartHandler){
	cart:=api.Group("/cart")
	routes:=cart.Group("/")
	routes.Use(middleware.UserMiddleware)
	{
		routes.Post("/add",cartH.AddToCart)
		routes.Delete("/clear", cartH.ClearCart)
		routes.Delete("/:id", cartH.RemoveItem)
		routes.Get("/",cartH.GetCart)
		routes.Put("/:id", cartH.UpdateQuantity)
	}

}