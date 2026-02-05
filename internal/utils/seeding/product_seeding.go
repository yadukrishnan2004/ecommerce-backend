package seeding

import (
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) {

	var count int64
	db.Model(&domain.Product{}).Count(&count)

	if count > 0 {
		fmt.Println("Products already exist. Skipping product seed.")
		return
	}

	products := []domain.Product{
		{
			Name:        "iPhone 15 Pro",
			Description: "The ultimate iPhone with titanium design.",
			Price:       999.00,
			Stock:       50,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/000/fff&text=iPhone+Front", "https://dummyimage.com/400x400/000/fff&text=iPhone+Back"},
		},
		{
			Name:        "MacBook Air M3",
			Description: "Supercharged by the M3 chip. Portable and powerful.",
			Price:       1299.00,
			Stock:       30,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/333/fff&text=MacBook"},
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Industry-leading noise canceling headphones.",
			Price:       349.99,
			Stock:       100,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/555/fff&text=Headphones"},
		},
		{
			Name:        "Nike Air Jordan 1",
			Description: "Classic high-top sneakers. Retro style.",
			Price:       180.00,
			Stock:       10,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/red/fff&text=Jordan+1"},
		},
		{
			Name:        "Mechanical Keyboard",
			Description: "RGB Backlit, Cherry MX Blue switches.",
			Price:       89.50,
			Stock:       25,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/blue/fff&text=Keyboard"},
		},
		{
			Name:        "4K Monitor 27-inch",
			Description: "Ultra HD resolution for professional editing.",
			Price:       450.00,
			Stock:       15,
			Images:      pq.StringArray{"https://dummyimage.com/400x400/111/fff&text=Monitor"},
		},
	}

	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	fmt.Println("Dummy products seeded successfully!")
}