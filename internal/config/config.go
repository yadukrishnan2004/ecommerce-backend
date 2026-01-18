package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App_Name string
	App_Port string
	DSN 	 string
}

type JWTConfig struct {
	Secret 		string
	AccessTTL 	time.Duration
	RefreshTTL  time.Duration
}


func getEnv(key, fallback string) string {
	if v:=os.Getenv(key);v != "" {
		return v
	}
	return fallback
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{

		App_Name: getEnv("App_name","ecommerce-backend"),
		App_Port: getEnv("APP_PORT","8080"),
		DSN: 	  getEnv("DSN",""),
	}
	return config
}

