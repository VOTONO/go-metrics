package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	address         string
	dbAdress        string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func getConfig() Config {
	config := Config{
		address:         "localhost:8080",
		dbAdress:        "localhost:5432",
		storeInterval:   300,
		fileStoragePath: "/tmp/metrics-db.json",
		restore:         true,
	}

	// Parse flags
	addressFlag := flag.String("a", config.address, "Address to bind to (default: localhost:8080)")
	dbAdressFlag := flag.String("d", config.dbAdress, "Address to bind db (default: localhost:5432)")
	storeIntervalFlag := flag.Int("i", config.storeInterval, "Store interval in seconds (default: 300)")
	fileStoragePathFlag := flag.String("f", config.fileStoragePath, "File storage path (default: /tmp/metrics-db.json)")
	restoreFlag := flag.Bool("r", config.restore, "Restore from file storage (default: true)")
	flag.Parse()

	// Override with command-line flags if provided
	config.address = *addressFlag
	config.dbAdress = *dbAdressFlag
	config.storeInterval = *storeIntervalFlag
	config.fileStoragePath = *fileStoragePathFlag
	config.restore = *restoreFlag

	// Override with environment variables if they exist
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.address = address
	}
	if dbAddress, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.dbAdress = dbAddress
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
