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
		Name:        "Sony WH-1000XM5",
		Description: "Industry-leading noise canceling headphones with up to 30 hours battery.",
		Price:       29999,
		Stock:       10,
		Category:    "wireless",
		Images:      pq.StringArray{"https://d1ncau8tqf99kp.cloudfront.net/converted/103364_original_local_1200x1050_v3_converted.webp"},
	},
	{
		Name:        "Boat Rockerz 550",
		Description: "Affordable wireless over-ear headphones with immersive sound.",
		Price:       1999,
		Stock:       50,
		Category:    "wireless",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/61gYxcIGjvL.SX522.jpg"},
	},
	{
		Name:        "JBL Tune 760NC",
		Description: "Lightweight, foldable headphones with active noise cancelling and deep bass.",
		Price:       7499,
		Stock:       25,
		Category:    "wireless",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/71TvdUf4kyL.SX522.jpg"},
	},
	{
		Name:        "Bose QuietComfort 45",
		Description: "Premium noise-cancelling headphones with balanced sound and all-day comfort.",
		Price:       29990,
		Stock:       10,
		Category:    "wireless",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/31n8zJ5mBDL.SX300_SY300_QL70_FMwebp.jpg"},
	},
	{
		Name:        "Sennheiser HD 450BT",
		Description: "Wireless headphones with deep dynamic bass, Bluetooth 5.0, and 30-hour battery.",
		Price:       8990,
		Stock:       20,
		Category:    "wireless",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/61h8NTXn5oL.UF1000,1000_QL80.jpg"},
	},
	{
		Name:        "Apple AirPods Max",
		Description: "Over-ear headphones with high-fidelity audio, computational sound, and spatial audio.",
		Price:       59900,
		Stock:       8,
		Category:    "luxury",
		Images:      pq.StringArray{"https://www.dimprice.co.uk/image/cache/png/apple-airpods-max/silver/apple-airpods-max-silver-04-550x550.png"},
	},
	{
		Name:        "Skullcandy Crusher Evo",
		Description: "Wireless headphones with adjustable sensory bass and 40 hours of battery life.",
		Price:       13999,
		Stock:       18,
		Category:    "wireless",
		Images:      pq.StringArray{"https://encrypted-tbn1.gstatic.com/shopping?q=tbn:ANd9GcSGhNDyBudvbHfqxY3zzRwxg6ZpGf_2TgFkHQ1BXhVmoyavvyHBT6nuvgtLsxVYa9_XEiPjJnfXSBKQ-wo7M9B6nYRbaIawS2g6YpHTTWcJ"},
	},
	{
		Name:        "OneOdio Pro-50",
		Description: "Wired studio headphones with 50mm drivers and professional sound quality.",
		Price:       3499,
		Stock:       30,
		Category:    "studio",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/71nT-y33zyL.SX679.jpg"},
	},
	{
		Name:        "AKG K371",
		Description: "Professional studio headphones tuned for accuracy with oval earcups.",
		Price:       11490,
		Stock:       12,
		Category:    "studio",
		Images:      pq.StringArray{"https://encrypted-tbn3.gstatic.com/shopping?q=tbn:ANd9GcSY_yAR5ezKIpPZD-JgmT8cviBpwCB1d-uzg0XEokP_gH95J0ZXKebXgaBxXBcKf689Ra4wiTdwgamuty7G5iB2tg8l9IwHrshaR9r4nVhDxj0V_hzWpOMTBg"},
	},
	{
		Name:        "Anker Soundcore Life Q30",
		Description: "Hybrid active noise cancelling headphones with 40 hours playtime and Hi-Res audio.",
		Price:       7999,
		Stock:       22,
		Category:    "wireless",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/31JvFCUXHZL.SX300_SY300_QL70_FMwebp.jpg"},
	},
	{
		Name:        "Sony WF-1000XM5",
		Description: "Flagship true wireless earbuds with industry-leading noise cancellation and immersive sound.",
		Price:       24990,
		Stock:       20,
		Category:    "earbuds",
		Images:      pq.StringArray{"https://store.sony.com.au/dw/image/v2/ABBC_PRD/on/demandware.static/-/Sites-sony-master-catalog/default/dw9194429e/images/WF1000XM5B/WF1000XM5B_2.png?sw=710&sh=710&sm=fit"},
	},
	{
		Name:        "boAt Rockerz 255 Pro+",
		Description: "Wireless neckband with 60 hours playback, IPX7 rating, and fast charging.",
		Price:       1099,
		Stock:       100,
		Category:    "neckband",
		Images:      pq.StringArray{"https://m.media-amazon.com/images/I/51Lqm2rKvHL.SX522.jpg"},
	},
	{
		Name:        "Sennheiser Momentum 4 Wireless",
		Description: "High-end wireless headphones with 60-hour battery life and adaptive noise cancellation.",
		Price:       34990,
		Stock:       10,
		Category:    "luxury",
		Images:      pq.StringArray{"https://www.hifi-regler.de/images_c/fm/products/sennheiser/sennheiser_momentum_ws.p1140x855.jpg"},
	},
	{
		Name:        "Marshall Major IV",
		Description: "Iconic on-ear headphones with 80+ hours wireless playtime and classic rock sound.",
		Price:       11999,
		Stock:       25,
		Category:    "wireless",
		Images:      pq.StringArray{"https://www.marshallheadphones.com/dw/image/v2/BCQL_PRD/on/demandware.static/-/Sites-zs-master-catalog/default/dw08944168/images/marshall/headphones/major-iv/large/pos-marshall-major-iv-black-01.png"},
	},
	{
		Name:        "Jabra Elite 85h",
		Description: "Smart noise-cancelling headphones with rain resistance and 36-hour battery.",
		Price:       16999,
		Stock:       18,
		Category:    "wireless",
		Images:      pq.StringArray{"https://www.jabra.com/-/media/Images/Products/Jabra-Elite-85h/Product/elite_85h_titanium_02.png?w=600"},
	},
	{
		Name:        "Shure AONIC 50",
		Description: "Studio-quality wireless headphones with adjustable noise cancellation and premium build.",
		Price:       34990,
		Stock:       10,
		Category:    "studio",
		Images:      pq.StringArray{"https://audio46.com/cdn/shop/files/AONIC50_Gen2_SBH50G2-BK_Front_Swivel_OnWhite_HRcopy_1200x1200.jpg?v=1693578640"},
	},
	{
		Name:        "Audio-Technica ATH-M50x",
		Description: "Legendary studio headphones with accurate sound and durable design.",
		Price:       37000,
		Stock:       20,
		Category:    "studio",
		Images:      pq.StringArray{"https://th.bing.com/th/id/R.a8585cba5aa0cbaf9e0226e9aa105442"},
	},
	{
		Name:        "Samsung Galaxy Buds2 Pro",
		Description: "Hi-Fi wireless earbuds with 24-bit audio and intelligent ANC.",
		Price:       17999,
		Stock:       30,
		Category:    "earbuds",
		Images:      pq.StringArray{"https://pisces.bbystatic.com/image2/BestBuy_US/images/products/6510/6510544cv12d.jpg"},
	},
	{
		Name:        "Realme Buds Wireless 3",
		Description: "Affordable neckband with ANC, 40-hour battery, and fast charging.",
		Price:       1799,
		Stock:       60,
		Category:    "neckband",
		Images:      pq.StringArray{"https://adminapi.applegadgetsbd.com/storage/media/large/realme-Buds-Wireless-3-Bass-Yellow-5884.jpg"},
	},
	{
		Name:        "HyperX Cloud Alpha",
		Description: "Gaming headphones with dual chamber drivers and detachable mic.",
		Price:       8990,
		Stock:       15,
		Category:    "gaming",
		Images:      pq.StringArray{"https://cdn.shopify.com/s/files/1/0564/3612/9997/products/hyperx_cloud_alpha_blackred_4_detachable_2048x2048.jpg"},
	},
}


	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	fmt.Println("Dummy products seeded successfully!")
}
