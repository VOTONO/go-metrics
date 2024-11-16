package helpers

import (
	"syscall"
	"testing"
	"time"
)

func TestCreateSystemStopChannel(t *testing.T) {
	t.Run("Test channel receives SIGINT signal", func(t *testing.T) {
		stopChannel := CreateSystemStopChannel()

		// Send a SIGINT signal
		go func() {
			stopChannel <- syscall.SIGINT
		}()

		// Select with timeout to avoid hanging test
		select {
		case sig := <-stopChannel:
			if sig != syscall.SIGINT {
				t.Errorf("Expected SIGINT, but got %v", sig)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for SIGINT signal")
		}
	})

	t.Run("Test channel receives SIGTERM signal", func(t *testing.T) {
		stopChannel := CreateSystemStopChannel()

		// Send a SIGTERM signal
		go func() {
			stopChannel <- syscall.SIGTERM
		}()

		select {
		case sig := <-stopChannel:
			if sig != syscall.SIGTERM {
				t.Errorf("Expected SIGTERM, but got %v", sig)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for SIGTERM signal")
		}
	})
}

func TestCreateTicker(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected time.Duration
	}{
		{"1 second ticker", 1, 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticker := CreateTicker(tt.seconds)
			defer ticker.Stop() // Clean up ticker to prevent resource leaks

			select {
			case <-ticker.C:
				// Test succeeded if the ticker channel ticks within expected time
			case <-time.After(tt.expected + (tt.expected / 2)): // Allow some buffer
				t.Errorf("Expected ticker interval around %v, but got timeout", tt.expected)
			}
		})
	}
}
