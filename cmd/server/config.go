package main

import (
	"flag"
	"os"
)

type Config struct {
	address string
}

func getConfig() Config {
	config := Config{}

	// Default values
	config.address = "localhost:8080"

	// Override with environment variables if they exist
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.address = address
	}

	// Parse flags
	addressFlag := flag.String("a", config.address, "Address to bind to (default: localhost:8080)")
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag

	return config
}
