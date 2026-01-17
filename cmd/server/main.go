package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repositery"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/infrastructure"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/service"
)

func main() {
	DB,err:=infrastructure.ConnectPostgres()
	if err != nil {
		log.Fatal("falile to connect with Database")
	}

	app:=fiber.New()
	userRepo:=repositery.NewUserRepo(DB)
	userSVC:=service.NewUserService(userRepo)
	UserHandler:=handler.NewUserHandler(userSVC)

	api:=app.Group("/")
	api.Post("register",UserHandler.Register)

	app.Listen(":8080")

}