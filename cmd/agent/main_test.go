package main

import (
	"net/http"
	"testing"
)

func TestRequestFor(t *testing.T) {
	tests := []struct {
		name          string
		adress        string
		metric        Metric
		expectedURL   string
		expectedError bool
	}{
		{
			name:          "Valid Counter Metric",
			adress:        "localhost:8080",
			metric:        Metric{Name: "PollCount", Type: "counter", Value: 0},
			expectedURL:   "http://localhost:8080/update/PollCount/counter/0",
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := requestFor(test.metric, test.adress)
			if (err != nil) != test.expectedError {
				t.Fatalf("expected error: %v, got: %v", test.expectedError, err)
			}
			if req != nil {
				if req.URL.String() != test.expectedURL {
					t.Errorf("expected URL: %s, got: %s", test.expectedURL, req.URL.String())
				}
				if req.Method != http.MethodPost {
					t.Errorf("expected method: POST, got: %s", req.Method)
				}
				if req.Header.Get("Content-Type") != "text/plain" {
					t.Errorf("expected Content-Type: text/plain, got: %s", req.Header.Get("Content-Type"))
				}
			}
		})
	}
}
