package initilizers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gofiber/fiber/v2"
	auth "github.com/yadukrishnan2004/ecommerce-backend/internal/Auth"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/notifications"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/usecase"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/router"
	"gorm.io/gorm"
)

func InitializeDependencies(
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
	orderRepo := repository.NewOrderRepo(db)
   	addressRepo:= repository.NewAddressRepo(db)
	jwtService := auth.NewJwtService(cfg.JWT)


	// Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo, notifier, *jwtService, orderRepo,productRepo)
	adminUseCase := usecase.NewAdminUseCase(userRepo, productRepo, orderRepo)
	cartService := usecase.NewCartService(cartRepo, productRepo)
	wishService := usecase.NewWishlistService(wishRepo, productRepo)
	orderService := usecase.NewOrderService(orderRepo, cartRepo, productRepo)
	addressUsecase:=usecase.NewAddressUsecase(addressRepo)


	// Handlers
	userHandler := handler.NewUserHandler(userUseCase)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	cartHandler := handler.NewCartHandler(cartService)
	wishHandler := handler.NewWishlistHandler(wishService)
	orderHandler := handler.NewOrderHandler(orderService)
	addressHandler:=handler.NewAddressHandler(addressUsecase)

	// Routes
	router.SetUpRouter(app, userHandler, adminHandler, cartHandler, wishHandler, orderHandler,addressHandler)
}

func StartServer(app *fiber.App, cfg *config.Config) {
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