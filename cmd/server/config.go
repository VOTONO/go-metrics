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
	defaultSecretKey       = ""
	defaultEnableHttps     = false
	defaultPublicKeyPath   = ""
	defaultPrivateKeyPath  = ""
)

type Config struct {
	address         string
	DSN             string
	storeInterval   int
	fileStoragePath string
	restore         bool
	secretKey       string
	enableHttps     bool
	publicKeyPath   string
	privateKeyPath  string
}

func loadEnvConfig(config *Config) {
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
	if secretKey, ok := os.LookupEnv("KEY"); ok {
		config.secretKey = secretKey
	}
	if enableHttps, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		if val, err := strconv.ParseBool(enableHttps); err == nil {
			config.enableHttps = val
		}
	}
	if publicKeyPath, ok := os.LookupEnv("PUBLIC_CRYPTO_KEY"); ok {
		config.publicKeyPath = publicKeyPath
	}
	if privateKeyPath, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		config.privateKeyPath = privateKeyPath
	}
}

func parseFlags(config *Config) {
	addressFlag := flag.String("a", config.address, fmt.Sprintf("Address to bind to (default: %s)", defaultAddress))
	dbAddressFlag := flag.String("d", config.DSN, fmt.Sprintf("Address to bind db (default: %s)", defaultDSN))
	storeIntervalFlag := flag.Int("i", config.storeInterval, fmt.Sprintf("Store interval in seconds (default: %d)", defaultStoreInterval))
	fileStoragePathFlag := flag.String("f", config.fileStoragePath, fmt.Sprintf("File storage path (default: %s)", defaultFileStoragePath))
	restoreFlag := flag.Bool("r", config.restore, fmt.Sprintf("Restore from file storage (default: %t)", defaultRestore))
	secretKeyFlag := flag.String("k", config.secretKey, fmt.Sprintf("Secret key (default: %s)", defaultSecretKey))
	enableHttpsFlag := flag.Bool("s", config.enableHttps, fmt.Sprintf("Enable HTTPS support (default: %t)", defaultEnableHttps))
	publicKeyPathFlag := flag.String("public-crypto-key", config.publicKeyPath, fmt.Sprintf("Public key path (default: %s)", defaultPublicKeyPath))
	privateKeyPathFlag := flag.String("crypto-key", config.privateKeyPath, fmt.Sprintf("Private key path (default: %s)", defaultPrivateKeyPath))

	flag.Parse()

	config.address = *addressFlag
	config.DSN = *dbAddressFlag
	config.storeInterval = *storeIntervalFlag
	config.fileStoragePath = *fileStoragePathFlag
	config.restore = *restoreFlag
	config.secretKey = *secretKeyFlag
	config.enableHttps = *enableHttpsFlag
	config.publicKeyPath = *publicKeyPathFlag
	config.privateKeyPath = *privateKeyPathFlag
}

func getConfig() Config {
	config := Config{
		address:         defaultAddress,
		DSN:             defaultDSN,
		storeInterval:   defaultStoreInterval,
		fileStoragePath: defaultFileStoragePath,
		restore:         defaultRestore,
		secretKey:       defaultSecretKey,
		enableHttps:     defaultEnableHttps,
		publicKeyPath:   defaultPublicKeyPath,
		privateKeyPath:  defaultPrivateKeyPath,
	}

	loadEnvConfig(&config)
	parseFlags(&config)

	return config
}
