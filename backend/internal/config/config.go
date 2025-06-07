package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath      string
	JWTSecret   string
	ServerPort  string
	NATSUrl     string
	Environment string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &Config{
		DBPath:      getEnv("DB_PATH", "data.db"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		NATSUrl:     getEnv("NATS_URL", "nats://localhost:4222"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 