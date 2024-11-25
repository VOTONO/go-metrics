package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultSecretKey      = ""
	defaultRateLimit      = 3
	defaultPublicKeyPath  = ""
	defaultConfigFilePath = ""
)

type Config struct {
	Address        string
	PollInterval   int
	ReportInterval int
	SecretKey      string
	RateLimit      int
	PublicKeyPath  string
}

func parseConfigFile(config *Config) {
	// Define a flag for the configuration file path
	configFilePath := flag.String("c", defaultConfigFilePath, fmt.Sprintf("Configuration file path (default: %s)", defaultConfigFilePath))

	// Parse the flags to get the file path if specified
	flag.Parse()

	// If no config file is specified, return without doing anything
	if *configFilePath == "" {
		return
	}

	// Open the configuration file
	file, err := os.Open(*configFilePath)
	if err != nil {
		log.Printf("Unable to open configuration file: %v", err)
		return
	}
	defer file.Close()

	// Read the file contents
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Unable to read file stats: %v", err)
		return
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		log.Printf("Error reading configuration file: %v", err)
		return
	}

	// Unmarshal JSON directly into the Config struct
	err = json.Unmarshal(buffer, config)
	if err != nil {
		log.Printf("Error decoding configuration file: %v", err)
		return
	}
}

func parseEnvs(config *Config) {
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.Address = address
	}
	if pollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if val, err := strconv.Atoi(pollInterval); err == nil {
			config.PollInterval = val
		}
	}
	if reportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if val, err := strconv.Atoi(reportInterval); err == nil {
			config.ReportInterval = val
		}
	}
	if secretKey, ok := os.LookupEnv("KEY"); ok {
		config.SecretKey = secretKey
	}
	if rateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		if val, err := strconv.Atoi(rateLimit); err == nil {
			config.RateLimit = val
		}
	}
	if publicKeyPath, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		config.PublicKeyPath = publicKeyPath
	}
}

func parseFlags(config *Config) {
	addressFlag := flag.String("a", config.Address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	pollIntervalFlag := flag.Int("p", config.PollInterval, fmt.Sprintf("Poll interval in seconds (default: %d)", defaultPollInterval))
	reportIntervalFlag := flag.Int("r", config.ReportInterval, fmt.Sprintf("Report interval in seconds (default: %d)", defaultReportInterval))
	secretKeyFlag := flag.String("k", config.SecretKey, fmt.Sprintf("Secret key (default: %s)", defaultSecretKey))
	rateLimitFlag := flag.Int("l", config.RateLimit, fmt.Sprintf("Rate limit key (default: %d)", defaultRateLimit))
	publicKeyPath := flag.String("crypto-key", config.PublicKeyPath, fmt.Sprintf("Public key path (default: %s)", defaultPublicKeyPath))

	flag.Parse()

	config.Address = *addressFlag
	config.PollInterval = *pollIntervalFlag
	config.ReportInterval = *reportIntervalFlag
	config.SecretKey = *secretKeyFlag
	config.RateLimit = *rateLimitFlag
	config.PublicKeyPath = *publicKeyPath
}

func getConfig() Config {
	config := Config{
		Address:        defaultAddress,
		PollInterval:   defaultPollInterval,
		ReportInterval: defaultReportInterval,
		SecretKey:      defaultSecretKey,
		RateLimit:      defaultRateLimit,
		PublicKeyPath:  defaultPublicKeyPath,
	}

	parseConfigFile(&config)
	parseEnvs(&config)
	parseFlags(&config)

	return config
}
