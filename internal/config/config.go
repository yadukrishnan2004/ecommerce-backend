package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	Secret     string
	AccessTTL  int32
	RefreshTTL int32
}
type Config struct {
	App_Name       string
	App_Port       string
	DSN            string
	JWT            *JWTConfig
	SMTP_EMAIL     string
	SMTP_PASS      string
	SMTP_HOST      string
	SMTP_PORT      string
	SMTP_FROM      string
	ADMIN_EMAIL    string
	ADMIN_PASSWORD string

	RAZORPAY_KEY    string
	RAZORPAY_SECRET string

	AllowedOrigins string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	jwtConfig := &JWTConfig{
		Secret:     getEnv("JWT_SECRET", ""),
		AccessTTL:  15 * 60,
		RefreshTTL: 7 * 24 * 60 * 60,
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	config := &Config{
		App_Name:       getEnv("App_name", "ecommerce-backend"),
		App_Port:       port,
		DSN:            getEnv("DSN", ""),
		JWT:            jwtConfig,
		SMTP_EMAIL:     getEnv("SMTP_EMAIL", ""),
		SMTP_PASS:      getEnv("SMTP_PASS", ""),
		SMTP_HOST:      getEnv("SMTP_HOST", ""),
		SMTP_PORT:      getEnv("SMTP_PORT", ""),
		SMTP_FROM:      getEnv("SMTP_FROM", ""),
		ADMIN_EMAIL:    getEnv("ADMIN_EMAIL", ""),
		ADMIN_PASSWORD: getEnv("ADMIN_PASSWORD", ""),

		RAZORPAY_KEY:    getEnv("RAZORPAY_KEY", ""),
		RAZORPAY_SECRET: getEnv("RAZORPAY_SECRET", ""),

		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173"),
	}
	return config
}
