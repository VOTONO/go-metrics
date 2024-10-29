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
