package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	auth "github.com/yadukrishnan2004/ecommerce-backend/internal/Auth"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/notifications"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/infrastructure"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/router"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/usecase"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Database Connection
	db := connectDB(cfg)

	// 3. Initialize Fiber App
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// 4. Initialize Dependencies & Routes
	initializeDependencies(app, db, cfg)

	// 5. Start Server with Graceful Shutdown
	startServer(app, cfg)
}

func connectDB(cfg *config.Config) *gorm.DB {
	db := infrastructure.ConnectPostgres(cfg.DSN)
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB connection: %v", err)
	}

	// Verify connection
	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(
		&repository.User{},
		&repository.Product{},
		&repository.CartItem{},
		&repository.Wishlist{},
		); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	return db
}

func initializeDependencies(
	app *fiber.App, 
	db *gorm.DB,
	cfg *config.Config,
	) {
	// Services & Adapters
	notifier := notifications.NewEmailNotifier(
		cfg.SMTP_HOST,
		cfg.SMTP_PORT,
		cfg.SMTP_EMAIL,
		cfg.SMTP_PASS,
	)

	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepo(db)
	cartRepo := repository.NewCartRepo(db)
	wishRepo := repository.NewWishlistRepo(db)
	jwtService := auth.NewJwtService(cfg.JWT)

	// Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo, notifier, *jwtService)
	adminUseCase := usecase.NewAdminUseCase(userRepo,productRepo)
	cartService := usecase.NewCartService(cartRepo, productRepo)
	wishService := usecase.NewWishlistService(wishRepo, productRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userUseCase)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	cartHandler := handler.NewCartHandler(cartService)
	wishHandler := handler.NewWishlistHandler(wishService)

	// Routes
	router.SetUpRouter(app, userHandler, adminHandler,cartHandler,wishHandler)
}

func startServer(app *fiber.App, cfg *config.Config) {
	// Run server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.App_Port)
		fmt.Printf("Server is running on %s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Panicf("Server error: %v", err)
		}
	}()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block until signal received

	fmt.Println("\nGracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server shutdown successfully")
}
