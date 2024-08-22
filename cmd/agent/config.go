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
	defaltSecretKey       = ""
)

type Config struct {
	address        string
	pollInterval   int
	reportInterval int
	secretKey      string
}

func getConfig() Config {
	config := Config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		secretKey:      defaltSecretKey,
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

	if secretKey, ok := os.LookupEnv("KEY"); ok {
		config.secretKey = secretKey
	}

	// Parse flags
	addressFlag := flag.String("a", config.address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	pollIntervalFlag := flag.Int("p", config.pollInterval, fmt.Sprintf("Poll interval in seconds (default: %d)", defaultPollInterval))
	reportIntervalFlag := flag.Int("r", config.reportInterval, fmt.Sprintf("Report interval in seconds (default: %d)", defaultReportInterval))
	secretKeyFlag := flag.String("k", config.secretKey, fmt.Sprintf("Secret key (default: %s)", defaltSecretKey))
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag
	config.pollInterval = *pollIntervalFlag
	config.reportInterval = *reportIntervalFlag
	config.secretKey = *secretKeyFlag

	return config
}
