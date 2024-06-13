package storage_test

import (
	"testing"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/storage"
)

func TestStorage(t *testing.T) {
	tests := []struct {
		name           string
		initialStorage map[string]models.Metric
		metricToStore  models.Metric
		metricToGet    models.Metric
	}{
		{
			name:           "Store gauge on empty storage",
			initialStorage: make(map[string]models.Metric),
			metricToStore:  models.Metric{Name: "foo", Type: "gauge", Value: "100"},
			metricToGet:    models.Metric{Name: "foo", Type: "gauge", Value: "100"},
		},
		{
			name:           "Store counter on empty storage",
			initialStorage: make(map[string]models.Metric),
			metricToStore:  models.Metric{Name: "foo", Type: "counter", Value: "100"},
			metricToGet:    models.Metric{Name: "foo", Type: "counter", Value: "100"},
		},
		{
			name:           "Replace existing",
			initialStorage: map[string]models.Metric{"foo": {Name: "foo", Type: "gauge", Value: "100"}},
			metricToStore:  models.Metric{Name: "foo", Type: "gauge", Value: "100"},
			metricToGet:    models.Metric{Name: "foo", Type: "gauge", Value: "100"},
		},
		{
			name:           "Store counter type with existing value int64",
			initialStorage: map[string]models.Metric{"foo": {Name: "foo", Type: "counter", Value: "100"}},
			metricToStore:  models.Metric{Name: "foo", Type: "counter", Value: "100"},
			metricToGet:    models.Metric{Name: "foo", Type: "counter", Value: "200"},
		},
		{
			name:           "Store counter type with existing value float64",
			initialStorage: map[string]models.Metric{"foo": {Name: "foo", Type: "counter", Value: "1.5"}},
			metricToStore:  models.Metric{Name: "foo", Type: "counter", Value: "1.5"},
			metricToGet:    models.Metric{Name: "foo", Type: "counter", Value: "3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stor := storage.New(test.initialStorage)

			err := stor.Store(test.metricToStore)
			if err != nil {
				t.Fatalf("returned an unexpected error: %v", err)
			}

			metric, exist := stor.Get(test.metricToGet.Name)

			if !exist {
				t.Errorf("expected metric not found")
			}

			if metric != test.metricToGet {
				t.Errorf("Expected value to be %v, got %v", test.metricToGet, metric)
			}
		})
	}
}
