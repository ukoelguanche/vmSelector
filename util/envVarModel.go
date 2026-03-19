package util

import (
	"log"
	"os"
	"strconv"
)

func Getenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Env var [%s] not set", key)
	}
	return value
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Env var [%s] must be an integer", key)
	}

	return intValue
}
