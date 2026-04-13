package config

import (
	"os"
)

// Config stores application configuration
type Config struct {
	DBURL   string
	AppPort string
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		DBURL:   getEnv("DATABASE_URL", "host=postgres user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable"),
		AppPort: getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
