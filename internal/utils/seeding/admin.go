package seeding

import (
	"fmt"
	"log"

	"github.com/yadukrishnan2004/ecommerce-backend/helper"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/config"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
	"gorm.io/gorm"
)

func AdminSeeding(db *gorm.DB,cfg *config.Config){

	var admin domain.User
	email:=cfg.ADMIN_EMAIL
	password:=cfg.ADMIN_PASSWORD
	if email == "" || password == "" {
		log.Println("Warning: ADMIN_EMAIL or ADMIN_PASSWORD not set in .env. Skipping admin seed.")
		return
	}

	if err := db.Where("email=?",email).First(&admin).Error; err == nil{
		fmt.Println("Admin already exists. Skipping seed.")
		return
	}

	hash,_:=helper.Hash(password)
	newAdmin:=domain.User{
		Name: "Admin",
		Email: email,
		Password: hash,
		Role: "admin",
		IsActive: true,
	}
	if err := db.Create(&newAdmin).Error; err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}
	fmt.Println("Admin user created successfully!")
}