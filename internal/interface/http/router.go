package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/infrasturcture/memory"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/interface/http/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase/user"
)

func RegisterRouter(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	userRepo := memory.NewUserRepo()
	registerUC := user.NewRegisterUseCase(userRepo)
	loginUserUC:=user.NewLoginUseCase(userRepo)

	userHandler := handler.NewUserHandler(registerUC,loginUserUC)

	app.Post("/users/register", userHandler.Register)
	app.Post("/users/login", userHandler.Login)
}
