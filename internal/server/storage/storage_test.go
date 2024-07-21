package storage_test

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/storage"
	"go.uber.org/zap"
)

func float64Ptr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64       { return &v }

func TestStorage(t *testing.T) {
	tests := []struct {
		name           string
		initialStorage map[string]models.Metric
		metricToStore  models.Metric
		expectedMetric models.Metric
	}{
		{
			name:           "Store gauge on empty storage",
			initialStorage: make(map[string]models.Metric),
			metricToStore:  models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(100)},
			expectedMetric: models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(100)},
		},
		{
			name:           "Store counter on empty storage",
			initialStorage: make(map[string]models.Metric),
			metricToStore:  models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(100)},
			expectedMetric: models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(100)},
		},
		{
			name:           "Replace existing gauge",
			initialStorage: map[string]models.Metric{"foo": {ID: "foo", MType: "gauge", Value: float64Ptr(100)}},
			metricToStore:  models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(200)},
			expectedMetric: models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(200)},
		},
		{
			name:           "Store counter type with existing value int64",
			initialStorage: map[string]models.Metric{"foo": {ID: "foo", MType: "counter", Delta: int64Ptr(100)}},
			metricToStore:  models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(100)},
			expectedMetric: models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(200)},
		},
		{
			name:           "Store counter type with existing value float64",
			initialStorage: map[string]models.Metric{"foo": {ID: "foo", MType: "counter", Delta: int64Ptr(150)}},
			metricToStore:  models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(150)},
			expectedMetric: models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(300)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Initialize zap logger
			logger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}
			defer logger.Sync()
			zapLogger := *logger.Sugar()

			// Initialize sqlmock
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			// Create storage with mocked db
			stor := storage.New(test.initialStorage, db, zapLogger)

			_, err = stor.Store(test.metricToStore)
			if err != nil {
				t.Fatalf("returned an unexpected error: %v", err)
			}

			metric, exist := stor.Get(test.metricToStore.ID)
			if !exist {
				t.Errorf("expected metric not found")
			}

			if !compareMetrics(metric, test.expectedMetric) {
				t.Errorf("expected value to be %v, got %v", test.expectedMetric, metric)
			}
		})
	}
}

func compareMetrics(a, b models.Metric) bool {
	if a.ID != b.ID || a.MType != b.MType {
		return false
	}

	if a.Delta != nil && b.Delta != nil && *a.Delta != *b.Delta {
		return false
	} else if (a.Delta != nil && b.Delta == nil) || (a.Delta == nil && b.Delta != nil) {
		return false
	}

	if a.Value != nil && b.Value != nil && *a.Value != *b.Value {
		return false
	} else if (a.Value != nil && b.Value == nil) || (a.Value == nil && b.Value != nil) {
		return false
	}

	return true
}
