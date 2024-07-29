package repo_test

import (
	"database/sql"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"go.uber.org/zap"
	"log"
	"testing"
)

func float64Ptr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64       { return &v }

// Helper function to create a new MetricStorerImpl instance with the given initial storage.
func newMetricStorerImpl(initialMetrics map[string]models.Metric) *repo.MetricStorerImpl {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	zapLogger := *logger.Sugar()

	return repo.New(false, "", 0, &sql.DB{}, &zapLogger)
}

func TestStoreGauge(t *testing.T) {
	stor := newMetricStorerImpl(make(map[string]models.Metric))
	metric := models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(100)}

	_, err := stor.Store(metric)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	storedMetric, exists := stor.Get(metric.ID)
	if !exists {
		t.Errorf("expected metric not found")
	}

	if !compareMetrics(storedMetric, metric) {
		t.Errorf("Expected value to be %v, got %v", metric, storedMetric)
	}
}

func TestStoreCounter(t *testing.T) {
	stor := newMetricStorerImpl(make(map[string]models.Metric))
	metric := models.Metric{ID: "foo", MType: "counter", Delta: int64Ptr(100)}

	_, err := stor.Store(metric)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	storedMetric, exists := stor.Get(metric.ID)
	if !exists {
		t.Errorf("expected metric not found")
	}

	if !compareMetrics(storedMetric, metric) {
		t.Errorf("Expected value to be %v, got %v", metric, storedMetric)
	}
}

func TestReplaceGauge(t *testing.T) {
	initialMetrics := map[string]models.Metric{
		"foo": {ID: "foo", MType: "gauge", Value: float64Ptr(100)},
	}
	stor := newMetricStorerImpl(initialMetrics)
	metric := models.Metric{ID: "foo", MType: "gauge", Value: float64Ptr(200)}

	_, err := stor.Store(metric)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	storedMetric, exists := stor.Get(metric.ID)
	if !exists {
		t.Errorf("expected metric not found")
	}

	if !compareMetrics(storedMetric, metric) {
		t.Errorf("Expected value to be %v, got %v", metric, storedMetric)
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
