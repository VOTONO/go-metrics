package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	address        string
	pollInterval   int
	reportInterval int
}

func getConfig() Config {
	config := Config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
	}

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
	addressFlag := flag.String("a", config.address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	pollIntervalFlag := flag.Int("p", config.pollInterval, fmt.Sprintf("Poll interval in seconds (default: %d)", defaultPollInterval))
	reportIntervalFlag := flag.Int("r", config.reportInterval, fmt.Sprintf("Report interval in seconds (default: %d)", defaultReportInterval))
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag
	config.pollInterval = *pollIntervalFlag
	config.reportInterval = *reportIntervalFlag

	return config
}
