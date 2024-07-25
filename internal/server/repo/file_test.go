package repo

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/VOTONO/go-metrics/internal/models"

	"go.uber.org/zap/zaptest"
)

func createTempFile(t *testing.T, content []byte) *os.File {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if content != nil {
		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
		if err := tmpFile.Sync(); err != nil {
			t.Fatalf("failed to sync temp file: %v", err)
		}
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile
}

func TestRead(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	metrics := map[string]models.Metric{
		"metric1": {
			ID:    "metric1",
			MType: "gauge",
			Value: func(f float64) *float64 { return &f }(123.45),
		},
		"metric2": {
			ID:    "metric2",
			MType: "counter",
			Delta: func(i int64) *int64 { return &i }(678),
		},
	}

	content, err := json.Marshal(metrics)
	if err != nil {
		t.Fatalf("failed to marshal metrics: %v", err)
	}

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		args    string
		want    map[string]models.Metric
		wantErr bool
	}{
		{
			name:    "Read existing file",
			args:    tmpFile.Name(),
			want:    metrics,
			wantErr: false,
		},
		{
			name:    "Read non-existent file",
			args:    "non_existent_file.json",
			want:    map[string]models.Metric{},
			wantErr: false,
		},
		{
			name:    "Read invalid file",
			args:    createTempFile(t, []byte("invalid content")).Name(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	metrics := map[string]models.Metric{
		"metric1": {
			ID:    "metric1",
			MType: "gauge",
			Value: func(f float64) *float64 { return &f }(123.45),
		},
		"metric2": {
			ID:    "metric2",
			MType: "counter",
			Delta: func(i int64) *int64 { return &i }(678),
		},
	}

	tmpFile := createTempFile(t, nil)
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		args    string
		metrics map[string]models.Metric
		wantErr bool
	}{
		{
			name:    "Write valid metrics",
			args:    tmpFile.Name(),
			metrics: metrics,
			wantErr: false,
		},
		{
			name:    "Write to invalid file",
			args:    "/invalid_path/metrics.json",
			metrics: metrics,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Write(tt.args, tt.metrics, logger); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, err := Read(tt.args, logger)
				if err != nil {
					t.Errorf("Read() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.metrics) {
					t.Errorf("Write() got = %v, want %v", got, tt.metrics)
				}
			}
		})
	}
}
