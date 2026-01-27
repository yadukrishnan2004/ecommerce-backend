package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
)

func SetUpRouther(app *fiber.App, userH *handler.UserHandler,AdminH *handler.AdminHandler) {
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	userRoutes := v1.Group("/users")

	// PUBLIC ROUTES (No Token Required)
	{
		userRoutes.Post("/signup", userH.Register)
		userRoutes.Post("/login", userH.Login)
		userRoutes.Post("/forgot-password", userH.Forgetpassword)
	}

	// SPECIAL ROUTES
	{
		userRoutes.Post("/verify",middleware.ResetMiddleware,userH.OtpVerify)
		userRoutes.Post("/reset-password", middleware.ResetMiddleware, userH.Resetpassword)
	}

	user := userRoutes.Group("/")
	user.Use(middleware.UserMiddleware)

	{
		user.Post("/logout", userH.Logout)
		user.Put("/profile", userH.UpdateProfile)
		user.Get("/profile", userH.GetProfile)
	}

}
