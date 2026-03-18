package util

import (
	"log"

	"github.com/joho/godotenv"
)

type Context struct {
	DevMode       string
	GpuString     string
	CentineVMName string
	Port          string
	PveHosst      string
	PveTokenId    string
	PveSecret     string
}

var ContextStorage = &Context{}

func LoadContext() {
	godotenv.Load()

	ContextStorage = &Context{
		GetEnvOrDefault("DEV_MODE", "false"),
		Getenv("GPU_STRING"),
		Getenv("CENTINEL_VM_NAME"),
		Getenv("PORT"),
		Getenv("PVE_HOST"),
		Getenv("PVE_TOKEN_ID"),
		Getenv("PVE_SECRET"),
	}
	log.Printf("ContextStorage loaded successfully\n")
}
