package seeding

import (
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/repository"
	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) {
	var count int64
	// Use the struct directly for counting
	db.Model(&repository.Product{}).Count(&count)

	if count > 0 {
		fmt.Println("Products already exist. Skipping product seed.")
		return
	}

	products := []repository.Product{
		{
			Name:        "iPhone 15 Pro",
			Description: "The ultimate iPhone with titanium design.",
			Price:       999.00,
			Stock:       50,
			Category:    "game",
			Images:      pq.StringArray{"https://dummyimage.com/front.jpg", "https://dummyimage.com/back.jpg"},
		},
		{
			Name:        "MacBook Air M3",
			Description: "Supercharged by the M3 chip.",
			Price:       1299.00,
			Stock:       30,
			Category:    "studio",
			Images:      pq.StringArray{"https://dummyimage.com/macbook.jpg"},
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Noise canceling headphones.",
			Price:       349.99,
			Stock:       100,
			Category:    "game",
			Images:      pq.StringArray{"https://dummyimage.com/headphones.jpg"},
		},
		{
			Name:        "Nike Air Jordan 1",
			Description: "Classic high-top sneakers.",
			Price:       180.00,
			Stock:       10,
			Category:    "casual",
			Images:      pq.StringArray{"https://dummyimage.com/jordan.jpg"},
		},
		{
			Name:        "Mechanical Keyboard",
			Description: "RGB Backlit.",
			Price:       89.50,
			Stock:       25,
			Category:    "game",
			Images:      pq.StringArray{"https://dummyimage.com/keyboard.jpg"},
		},
		{
			Name:        "4K Monitor",
			Description: "Ultra HD resolution.",
			Price:       450.00,
			Stock:       15,
			Category:    "game",
			Images:      pq.StringArray{"https://dummyimage.com/monitor.jpg"},
		},
	}

	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	fmt.Println("Dummy products seeded successfully!")
}
