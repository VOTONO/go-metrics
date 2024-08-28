package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultAddress         = "localhost:8080"
	defaultDSN             = ""
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true
)

type Config struct {
	address         string
	DSN             string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func getConfig() Config {
	config := Config{
		address:         defaultAddress,
		DSN:             defaultDSN,
		storeInterval:   defaultStoreInterval,
		fileStoragePath: defaultFileStoragePath,
		restore:         defaultRestore,
	}

	// Parse flags
	addressFlag := flag.String("a", config.address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	dbAddressFlag := flag.String("d", config.DSN, fmt.Sprintf("Address to bind db (default: %s)", defaultDSN))
	storeIntervalFlag := flag.Int("i", config.storeInterval, fmt.Sprintf("StoreSingle interval in seconds (default: %d)", defaultStoreInterval))
	fileStoragePathFlag := flag.String("f", config.fileStoragePath, fmt.Sprintf("File storage path (default: %s)", defaultFileStoragePath))
	restoreFlag := flag.Bool("r", config.restore, fmt.Sprintf("Restore from file storage (default: %t)", defaultRestore))
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag
	config.DSN = *dbAddressFlag
	config.storeInterval = *storeIntervalFlag
	config.fileStoragePath = *fileStoragePathFlag
	config.restore = *restoreFlag

	// Override with environment variables if they exist
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.address = address
	}
	if dbAddress, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DSN = dbAddress
	}
	if storeInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if i, err := strconv.Atoi(storeInterval); err == nil {
			config.storeInterval = i
		}
	}
	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.fileStoragePath = fileStoragePath
	}
	if restore, ok := os.LookupEnv("RESTORE"); ok {
		if b, err := strconv.ParseBool(restore); err == nil {
			config.restore = b
		}
	}

	return config
}
