package main

import (
	"fmt"

	"github.com/razorpay/razorpay-go"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Loaded KEY: '%s'\n", cfg.RAZORPAY_KEY)
	fmt.Printf("Loaded SECRET: '%s'\n", cfg.RAZORPAY_SECRET)

	client := razorpay.NewClient(cfg.RAZORPAY_KEY, cfg.RAZORPAY_SECRET)
	data := map[string]interface{}{"amount": 50000, "currency": "INR", "receipt": "rx_1"}
	body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println("SUCCESS:", body["id"])
	}
}
