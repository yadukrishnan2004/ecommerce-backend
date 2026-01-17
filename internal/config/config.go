package config

import (
	"os"
	"time"
)

type Config struct {
	App_Name string
	App_Port string
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
	config := &Config{
		App_Name: getEnv("App_name","ecommerce-backend"),
		App_Port: getEnv("APP_PORT","8080"),
	}
	return config
}

