package cfg

import "os"

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" default:"postgres://items_user:items_pass@localhost:5433/items_db?sslmode=disable"`
}

func LoadConfig() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://items_user:items_pass@localhost:5433/items_db?sslmode=disable"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
