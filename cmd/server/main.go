package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/infrastructure"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/initilizers"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/seeding"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Database Connection
	db := infrastructure.ConnectDB(cfg)

	//admin seeding
	seeding.AdminSeeding(db, cfg)
	seeding.SeedProducts(db)

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", 
		AllowCredentials: true,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 4. Initialize Dependencies & Routes
	initilizers.InitializeDependencies(app, db, cfg)

	// 5. Start Server with Graceful Shutdown
	initilizers.StartServer(app, cfg)
}




