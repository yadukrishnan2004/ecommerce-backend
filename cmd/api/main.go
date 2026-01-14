package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/interface/http"
)

func main() {
	if err:=godotenv.Load();err !=nil{
		log.Println("no env file is found")
	} 
	cfg:=config.Load()
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// helper handler function 
	http.RegisterRouter(app)


		
	app.Listen(":"+cfg.App_Port)
}