package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestLoadEnvConfig(t *testing.T) {
	os.Setenv("ADDRESS", "127.0.0.1:9090")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("REPORT_INTERVAL", "15")
	os.Setenv("KEY", "secret123")
	os.Setenv("RATE_LIMIT", "7")

	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("REPORT_INTERVAL")
		os.Unsetenv("KEY")
		os.Unsetenv("RATE_LIMIT")
	}()

	config := Config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		secretKey:      defaultSecretKey,
		rateLimit:      defaultRateLimit,
	}
	loadEnvConfig(&config)

	expected := Config{
		address:        "127.0.0.1:9090",
		pollInterval:   5,
		reportInterval: 15,
		secretKey:      "secret123",
		rateLimit:      7,
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config with environment variables %v, got %v", expected, config)
	}
}

func TestParseFlags(t *testing.T) {
	// Reset flag.CommandLine to prevent interference
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Mock command-line arguments
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:9090",
		"-p", "5",
		"-r", "15",
		"-k", "secret123",
		"-l", "7",
	}

	config := Config{
		address:        defaultAddress,
		pollInterval:   defaultPollInterval,
		reportInterval: defaultReportInterval,
		secretKey:      defaultSecretKey,
		rateLimit:      defaultRateLimit,
	}
	parseFlags(&config)

	expected := Config{
		address:        "127.0.0.1:9090",
		pollInterval:   5,
		reportInterval: 15,
		secretKey:      "secret123",
		rateLimit:      7,
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config with flags %v, got %v", expected, config)
	}
}

func TestGetConfig(t *testing.T) {
	// Reset flag.CommandLine to prevent interference with flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set environment variables for the test
	os.Setenv("ADDRESS", "127.0.0.1:8081")
	os.Setenv("POLL_INTERVAL", "10")
	os.Setenv("REPORT_INTERVAL", "20")
	os.Setenv("KEY", "testkey")
	os.Setenv("RATE_LIMIT", "5")

	// Mock command-line arguments
	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:8082", // Override address from flag
		"-p", "15", // Override poll interval from flag
		"-r", "30", // Override report interval from flag
		"-k", "newsecret", // Override secret key from flag
		"-l", "10", // Override rate limit from flag
	}

	// Call the getConfig function to load configuration
	config := getConfig()

	// Expected values based on the environment variables and flags
	expected := Config{
		address:        "127.0.0.1:8082", // Flag overrides env
		pollInterval:   15,               // Flag overrides env
		reportInterval: 30,               // Flag overrides env
		secretKey:      "newsecret",      // Flag overrides env
		rateLimit:      10,               // Flag overrides env
	}

	// Clean up the environment variables after the test
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("REPORT_INTERVAL")
		os.Unsetenv("KEY")
		os.Unsetenv("RATE_LIMIT")
	}()

	// Check if the config matches the expected result
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected config: %+v, got: %+v", expected, config)
	}
}
