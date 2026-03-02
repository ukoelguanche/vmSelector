package util

import (
	"log"
	"strconv"
)

type Context struct {
	GpuString     string
	CentineVMName string
	PollInterval  int
	Port          string
	PveHosst      string
	PveTokenId    string
	PveSecret     string
}

var ContextStorage = &Context{}

func LoadContext() {
	pollInterval, err := strconv.Atoi(Getenv("POLL_INTERVAL"))
	if err != nil {
		log.Fatal("Could not parse POLL_INTERVAL value: %s\n", Getenv("POLL_INTERVAL"))
	}

	ContextStorage = &Context{

		Getenv("GPU_STRING"),
		Getenv("CENTINEL_VM_NAME"),
		pollInterval,
		Getenv("PORT"),
		Getenv("PVE_HOST"),
		Getenv("PVE_TOKEN_ID"),
		Getenv("PVE_SECRET"),
	}
	log.Printf("ContextStorage loaded successfully\n")
}
