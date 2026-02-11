package infrastructure

import (
	"log"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(dsn string) *gorm.DB {
	if dsn == "" {
		log.Fatal("no dsn found")
		return nil
	}
	Dsn := dsn
	DB, err := gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("faile to connect with the database")
		return nil
	}
	return DB
}

func ConnectDB(cfg *config.Config) *gorm.DB {
	db := ConnectPostgres(cfg.DSN)
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
		&repository.SignupRequest{},
		&repository.Product{},
		&repository.CartItem{},
		&repository.Wishlist{},
		&repository.Order{},
		&repository.OrderItem{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	return db
}
