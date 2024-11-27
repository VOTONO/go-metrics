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
		Address:         defaultAddress,
		DSN:             defaultDSN,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		SecretKey:       defaultSecretKey,
	}
	parseEnvs(&config)

	expected := Config{
		Address:         "127.0.0.1:9090",
		DSN:             "user:password@/dbname",
		StoreInterval:   600,
		FileStoragePath: "/custom/path/to/db.json",
		Restore:         false,
		SecretKey:       "my_secret_key",
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
		Address:         defaultAddress,
		DSN:             defaultDSN,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		SecretKey:       defaultSecretKey,
	}
	parseFlags(&config)

	expected := Config{
		Address:         "127.0.0.1:9090",
		DSN:             "user:password@/dbname",
		StoreInterval:   600,
		FileStoragePath: "/custom/path/to/db.json",
		Restore:         false,
		SecretKey:       "my_secret_key",
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config %v, got %v", expected, config)
	}
}

func TestGetConfig(t *testing.T) {
	// Reset flag.CommandLine to prevent interference with flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set environment variables for the test
	os.Setenv("ADDRESS", "127.0.0.1:8081")
	os.Setenv("DATABASE_DSN", "user:password@/mydb")
	os.Setenv("STORE_INTERVAL", "500")
	os.Setenv("FILE_STORAGE_PATH", "/new/path/to/db.json")
	os.Setenv("RESTORE", "true")
	os.Setenv("KEY", "new_secret_key")

	// Mock command-line arguments
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:8082", // Flag overrides environment variable for Address
		"-d", "user:password@/testdb", // Flag overrides environment variable for DSN
		"-i", "1000", // Flag overrides environment variable for store interval
		"-f", "/another/path/to/db.json", // Flag overrides environment variable for file storage path
		"-r=false",              // Flag overrides environment variable for Restore
		"-k", "override_secret", // Flag overrides environment variable for secret key
	}

	// Call the getConfig function to load configuration
	config := getConfig()

	// Expected config after combining environment variables and flags
	expected := Config{
		Address:         "127.0.0.1:8082",           // Flag overrides env
		DSN:             "user:password@/testdb",    // Flag overrides env
		StoreInterval:   1000,                       // Flag overrides env
		FileStoragePath: "/another/path/to/db.json", // Flag overrides env
		Restore:         false,                      // Flag overrides env
		SecretKey:       "override_secret",          // Flag overrides env
	}

	// Clean up the environment variables after the test
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("RESTORE")
		os.Unsetenv("KEY")
	}()

	// Check if the config matches the expected result
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config: %+v, got: %+v", expected, config)
	}
}
