package config

import "os"

type Config struct {
	App_Name string
	App_Port string
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