package cfg

import "os"

type Config struct {
	UserServiceURL string `env:"USER_SERVICE_URL" default:"http://localhost:8080"`
	ItemServiceURL string `env:"ITEM_SERVICE_URL" default:"http://localhost:8081"`
}

func LoadConfig() (*Config, error) {
	return &Config{
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8080"),
		ItemServiceURL: getEnv("ITEM_SERVICE_URL", "http://localhost:8081"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
