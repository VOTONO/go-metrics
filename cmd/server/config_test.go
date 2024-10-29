package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestLoadEnvConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("ADDRESS", "127.0.0.1:9090")
	os.Setenv("DATABASE_DSN", "user:password@/dbname")
	os.Setenv("STORE_INTERVAL", "600")
	os.Setenv("FILE_STORAGE_PATH", "/custom/path/to/db.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("KEY", "my_secret_key")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("ADDRESS")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("RESTORE")
		os.Unsetenv("KEY")
	}()

	config := Config{
		address:         defaultAddress,
		DSN:             defaultDSN,
		storeInterval:   defaultStoreInterval,
		fileStoragePath: defaultFileStoragePath,
		restore:         defaultRestore,
		secretKey:       defaultSecretKey,
	}
	loadEnvConfig(&config)

	expected := Config{
		address:         "127.0.0.1:9090",
		DSN:             "user:password@/dbname",
		storeInterval:   600,
		fileStoragePath: "/custom/path/to/db.json",
		restore:         false,
		secretKey:       "my_secret_key",
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config %v, got %v", expected, config)
	}
}

func TestParseFlags(t *testing.T) {
	// Reset flag.CommandLine to prevent interference with flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set command-line arguments
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:9090",
		"-d", "user:password@/dbname",
		"-i", "600",
		"-f", "/custom/path/to/db.json",
		"-r=false",
		"-k", "my_secret_key",
	}

	config := Config{
		address:         defaultAddress,
		DSN:             defaultDSN,
		storeInterval:   defaultStoreInterval,
		fileStoragePath: defaultFileStoragePath,
		restore:         defaultRestore,
		secretKey:       defaultSecretKey,
	}
	parseFlags(&config)

	expected := Config{
		address:         "127.0.0.1:9090",
		DSN:             "user:password@/dbname",
		storeInterval:   600,
		fileStoragePath: "/custom/path/to/db.json",
		restore:         false,
		secretKey:       "my_secret_key",
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config %v, got %v", expected, config)
	}
}
