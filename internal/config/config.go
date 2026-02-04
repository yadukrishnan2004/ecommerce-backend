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
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtConfig := &JWTConfig{
		Secret:     getEnv("JWT_SECRET", ""),
		AccessTTL:  15 * 60,
		RefreshTTL: 7 * 24 * 60 * 60,
	}

	config := &Config{
		App_Name:       getEnv("App_name", "ecommerce-backend"),
		App_Port:       getEnv("APP_PORT", "8080"),
		DSN:            getEnv("DSN", ""),
		JWT:            jwtConfig,
		SMTP_EMAIL:     getEnv("SMTP_EMAIL", ""),
		SMTP_PASS:      getEnv("SMTP_PASS", ""),
		SMTP_HOST:      getEnv("SMTP_HOST", ""),
		SMTP_PORT:      getEnv("SMTP_PORT", ""),
		SMTP_FROM:      getEnv("SMTP_FROM", ""),
		ADMIN_EMAIL:    getEnv("ADMIN_EMAIL", ""),
		ADMIN_PASSWORD: getEnv("ADMIN_PASSWORD", ""),
	}
	return config
}
