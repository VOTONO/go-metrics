package storage_test

import (
	"testing"

	"github.com/VOTONO/go-metrics/internal/server/storage"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		initialStorage map[string]interface{}
		keyToStore     string
		valueToStore   interface{}
		expectedValue  interface{}
	}{
		{
			name:           "Replace on empty storage",
			method:         "Replace",
			initialStorage: make(map[string]interface{}),
			keyToStore:     "key",
			valueToStore:   int64(100),
			expectedValue:  int64(100),
		},
		{
			name:           "Increment on empty storage",
			method:         "Increment",
			initialStorage: make(map[string]interface{}),
			keyToStore:     "key",
			valueToStore:   int64(100),
			expectedValue:  int64(100),
		},
		{
			name:           "Replace existing",
			method:         "Replace",
			initialStorage: map[string]interface{}{"key": int64(50)},
			keyToStore:     "key",
			valueToStore:   int64(100),
			expectedValue:  int64(100),
		},
		{
			name:           "Increment existing int64",
			method:         "Increment",
			initialStorage: map[string]interface{}{"key": int64(100)},
			keyToStore:     "key",
			valueToStore:   int64(100),
			expectedValue:  int64(200),
		},
		{
			name:           "Increment existing float64",
			method:         "Increment",
			initialStorage: map[string]interface{}{"key": float64(1.5)},
			keyToStore:     "key",
			valueToStore:   float64(1.5),
			expectedValue:  float64(3.0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memStorage := storage.New(test.initialStorage)

			switch test.method {
			case "Replace":
				err := memStorage.Replace(test.keyToStore, test.valueToStore)
				if err != nil {
					t.Fatalf("Replace returned an unexpected error: %v", err)
				}

				got := memStorage.Get(test.keyToStore)

				if got != test.expectedValue {
					t.Errorf("Expected value to be %v, got %v", test.expectedValue, got)
				}
			case "Increment":
				err := memStorage.Increment(test.keyToStore, test.valueToStore)
				if err != nil {
					t.Fatalf("Replace returned an unexpected error: %v", err)
				}

				got := memStorage.Get(test.keyToStore)

				if got != test.expectedValue {
					t.Errorf("Expected value to be %v, got %v", test.expectedValue, got)
				}
			}
		})
	}
}
