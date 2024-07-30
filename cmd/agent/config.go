package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	address        string
	pollInterval   int
	reportInterval int
}

func getConfig() Config {
	config := Config{}

	// Default values
	config.address = "localhost:8080"
	config.pollInterval = 2
	config.reportInterval = 10

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.address = address
	}

	if pollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if val, err := strconv.Atoi(pollInterval); err == nil {
			config.pollInterval = val
		} else {
			fmt.Printf("Invalid POLL_INTERVAL value: %v\n", pollInterval)
		}
	}

	if reportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if val, err := strconv.Atoi(reportInterval); err == nil {
			config.reportInterval = val
		} else {
			fmt.Printf("Invalid REPORT_INTERVAL value: %v\n", reportInterval)
		}
	}

	// Parse flags
	addressFlag := flag.String("a", config.address, "Address to bind to (default: localhost:8080)")
	pollIntervalFlag := flag.Int("p", config.pollInterval, "Poll interval in seconds (default: 2)")
	reportIntervalFlag := flag.Int("r", config.reportInterval, "Report interval in seconds (default: 10)")
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag
	config.pollInterval = *pollIntervalFlag
	config.reportInterval = *reportIntervalFlag

	return config
}
