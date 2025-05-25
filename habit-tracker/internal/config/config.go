package config

import (
	"os"
)

type Config struct {
	Port          string
	AWSRegion     string
	DynamoDBTable string
	Environment   string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		AWSRegion:     getEnv("AWS_REGION", "us-east-1"),
		DynamoDBTable: getEnv("DYNAMODB_TABLE", "habits"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
