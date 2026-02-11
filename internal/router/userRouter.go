package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetUpUserRouter(api fiber.Router, userH *handler.UserHandler) {

	userRoutes := api.Group("/users")

	// PUBLIC ROUTES
	{
		userRoutes.Post("/signup", userH.SignUp)
		userRoutes.Post("/verify", userH.OtpVerify)
		userRoutes.Post("/login", userH.SignIn)
		userRoutes.Post("/forgot-password", userH.Forgotpassword)
		userRoutes.Get("/allproducts", userH.GetAll)
		userRoutes.Get("/search", userH.SearchProducts)
		userRoutes.Get("/filter", userH.FilterProducts)
	}

	// SPECIAL ROUTES (Reset Token Required)
	{
		userRoutes.Post("/reset-password", middleware.ResetMiddleware, userH.Resetpassword)
	}

	// AUTHENTICATED ROUTES (JWT Required)

	protected := userRoutes.Group("/")
	protected.Use(middleware.UserMiddleware)

	{
		protected.Post("/logout", userH.Logout)
		protected.Put("/profile", userH.UpdateProfile)
		protected.Get("/profile", userH.GetProfile)
		protected.Get("/:id/orders", userH.GetOrder)
		protected.Put("/:id/cancel", userH.CancelOrder)
		protected.Get("/:id/orders/details", userH.GetOrderProduct)
	}
}
