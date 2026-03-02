package util

import (
	"log"
	"os"
)

func Getenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Env var [%s] not set", key)
	}
	return value
}
