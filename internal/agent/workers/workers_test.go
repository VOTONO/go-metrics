package workers

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func ExampleReadWorker_Start() {
	// Init logger for worker.
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	// Catch syscall to stop worker.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)
	defer stop()

	// Init worker.
	worker := NewReadWorker(
		sugaredLogger, // Metrics will be sent in this channel.
		2,             // Interval in seconds of metrics reading.
	)

	// Call Start() function in goroutine.
	go func() {
		worker.Start()
	}()

	select {
	// Metrics will be sent in this channel periodically.
	case readResult := <-worker.ResultChannel:
		fmt.Println(readResult)
	case <-ctx.Done():
		// Stop worker gracefully.
		worker.Stop()
	}
}

// TestNewReadWorker verifies the initialization of the ReadWorker struct
func TestNewReadWorker(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	interval := 5

	worker := NewReadWorker(logger, interval)

	if worker.logger == nil {
		t.Fatal("Logger should not be nil")
	}
	if worker.ticker == nil {
		t.Fatal("Ticker should not be nil")
	}
	if worker.stopChannel == nil {
		t.Fatal("Stop channel should not be nil")
	}
	if worker.ResultChannel == nil {
		t.Fatal("Result channel should not be nil")
	}
}

// TestReadMemoryMetrics verifies the memory metrics
func TestReadMemoryMetrics(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	worker := NewReadWorker(logger, 1)

	metrics := worker.readMemoryMetrics()

	// Verify key metrics are present
	expectedMetrics := map[string]bool{
		"TotalMemory": true,
		"FreeMemory":  true,
		"UsedMemory":  true,
	}
	for _, metric := range metrics {
		delete(expectedMetrics, metric.ID)
	}
	if len(expectedMetrics) > 0 {
		t.Errorf("Missing expected memory metrics: %v", expectedMetrics)
	}
}

// TestReadRuntimeMetrics verifies the runtime metrics
func TestReadRuntimeMetrics(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	worker := NewReadWorker(logger, 1)

	metrics := worker.readRuntimeMetrics()

	// Check for key metrics
	foundPollCount := false
	foundRandomValue := false
	for _, metric := range metrics {
		if metric.ID == "PollCount" {
			foundPollCount = true
		}
		if metric.ID == "RandomValue" {
			foundRandomValue = true
		}
	}
	if !foundPollCount || !foundRandomValue {
		t.Error("Missing key runtime metrics: PollCount or RandomValue")
	}
}

// TestFloat64Ptr verifies the float64Ptr helper function
func TestFloat64Ptr(t *testing.T) {
	val := 42.42
	ptr := float64Ptr(val)
	if ptr == nil || *ptr != val {
		t.Errorf("float64Ptr(%v) = %v, want %v", val, ptr, val)
	}
}
