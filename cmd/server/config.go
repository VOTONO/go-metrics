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
	defaultAddress         = "localhost:8080"
	defaultDSN             = ""
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true
	defaultSecretKey       = ""
	defaultEnableHTTPS     = false
	defaultPublicKeyPath   = ""
	defaultPrivateKeyPath  = ""
	defaultConfigFilePath  = ""
)

type Config struct {
	Address         string
	DSN             string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	SecretKey       string
	EnableHTTPS     bool
	PublicKeyPath   string
	PrivateKeyPath  string
}

func parseEnvs(config *Config) {
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		config.Address = address
	}
	if dbAddress, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DSN = dbAddress
	}
	if storeInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if i, err := strconv.Atoi(storeInterval); err == nil {
			config.StoreInterval = i
		}
	}
	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = fileStoragePath
	}
	if restore, ok := os.LookupEnv("RESTORE"); ok {
		if b, err := strconv.ParseBool(restore); err == nil {
			config.Restore = b
		}
	}
	if secretKey, ok := os.LookupEnv("KEY"); ok {
		config.SecretKey = secretKey
	}
	if enableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		if val, err := strconv.ParseBool(enableHTTPS); err == nil {
			config.EnableHTTPS = val
		}
	}
	if publicKeyPath, ok := os.LookupEnv("PUBLIC_CRYPTO_KEY"); ok {
		config.PublicKeyPath = publicKeyPath
	}
	if privateKeyPath, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		config.PrivateKeyPath = privateKeyPath
	}
}

func parseFlags(config *Config) {
	addressFlag := flag.String("a", config.Address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	dbAddressFlag := flag.String("d", config.DSN, fmt.Sprintf("Address to bind db (default: %s)", defaultDSN))
	storeIntervalFlag := flag.Int("i", config.StoreInterval, fmt.Sprintf("Store interval in seconds (default: %d)", defaultStoreInterval))
	fileStoragePathFlag := flag.String("f", config.FileStoragePath, fmt.Sprintf("File storage path (default: %s)", defaultFileStoragePath))
	restoreFlag := flag.Bool("r", config.Restore, fmt.Sprintf("Restore from file storage (default: %t)", defaultRestore))
	secretKeyFlag := flag.String("k", config.SecretKey, fmt.Sprintf("Secret key (default: %s)", defaultSecretKey))
	enableHTTPSFlag := flag.Bool("s", config.EnableHTTPS, fmt.Sprintf("Enable HTTPS support (default: %t)", defaultEnableHTTPS))
	publicKeyPathFlag := flag.String("public-crypto-key", config.PublicKeyPath, fmt.Sprintf("Public key path (default: %s)", defaultPublicKeyPath))
	privateKeyPathFlag := flag.String("crypto-key", config.PrivateKeyPath, fmt.Sprintf("Private key path (default: %s)", defaultPrivateKeyPath))

	flag.Parse()

	config.Address = *addressFlag
	config.DSN = *dbAddressFlag
	config.StoreInterval = *storeIntervalFlag
	config.FileStoragePath = *fileStoragePathFlag
	config.Restore = *restoreFlag
	config.SecretKey = *secretKeyFlag
	config.EnableHTTPS = *enableHTTPSFlag
	config.PublicKeyPath = *publicKeyPathFlag
	config.PrivateKeyPath = *privateKeyPathFlag
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

func getConfig() Config {
	config := Config{
		Address:         defaultAddress,
		DSN:             defaultDSN,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		SecretKey:       defaultSecretKey,
		EnableHTTPS:     defaultEnableHTTPS,
		PublicKeyPath:   defaultPublicKeyPath,
		PrivateKeyPath:  defaultPrivateKeyPath,
	}

	parseConfigFile(&config)
	parseEnvs(&config)
	parseFlags(&config)

	return config
}
