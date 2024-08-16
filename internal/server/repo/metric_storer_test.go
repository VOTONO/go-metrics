package repo_test

import (
	"context"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/handlers/utils"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

func TestMetricStorers(t *testing.T) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	zapLogger := *logger.Sugar()

	testMetricsFilePath := "/tmp/test_metrics.json"

	os.Remove(testMetricsFilePath)

	storers := []struct {
		name   string
		storer repo.MetricStorer
	}{
		{"LocalMetricStorer", repo.NewLocalMetricStorer(false, "", &zapLogger)},
		{"FileMetricStorer", repo.NewFileMetricStorer(testMetricsFilePath, &zapLogger)},
	}

	for _, stor := range storers {
		t.Run(stor.name, func(t *testing.T) {
			testStoreGetGauge(t, stor.storer)
			testStoreGetCounter(t, stor.storer)
			testInvalidMetric(t, stor.storer)
			testStoreSlice(t, stor.storer) // Add this line to include StoreSlice tests
		})
	}
}

func testStoreGetGauge(t *testing.T, stor repo.MetricStorer) {

	_, err := stor.StoreSingle(context.Background(), utils.ValidGaugeMetric)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	storedMetric, exists, getErr := stor.Get(context.Background(), utils.ValidGaugeMetric.ID)

	if getErr != nil {
		t.Fatalf("returned an unexpected error: %v", getErr)
	}

	if !exists {
		t.Errorf("expected metric not found")
	}

	if !compareMetrics(storedMetric, utils.ValidGaugeMetric) {
		t.Errorf("Expected value to be %v, got %v", utils.ValidGaugeMetric, storedMetric)
	}
}

func testStoreGetCounter(t *testing.T, stor repo.MetricStorer) {
	_, err := stor.StoreSingle(context.Background(), utils.ValidCounterMetric)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	storedMetric, exists, getErr := stor.Get(context.Background(), utils.ValidCounterMetric.ID)

	if getErr != nil {
		t.Fatalf("returned an unexpected error: %v", getErr)
	}

	if !exists {
		t.Errorf("expected metric not found")
	}

	if !compareMetrics(storedMetric, utils.ValidCounterMetric) {
		t.Errorf("Expected value to be %v, got %v", utils.ValidCounterMetric, storedMetric)
	}
}

func testInvalidMetric(t *testing.T, stor repo.MetricStorer) {
	invalidMetrics := []models.Metric{
		utils.InvalidGaugeMissingValue,
		utils.InvalidCounterMissingDelta,
		utils.InvalidMetricUnknownType,
		utils.InvalidMetricCounterWithGaugeValue,
		utils.InvalidMetricGaugeWithCounterDelta,
		utils.InvalidMetricEmptyID,
		utils.InvalidMetricNilBothFields,
	}

	for _, metric := range invalidMetrics {
		_, err := stor.StoreSingle(context.Background(), metric)
		if err == nil {
			t.Fatalf("expected an error for metric %v, but got none", metric)
		}
	}
}

func testStoreSlice(t *testing.T, stor repo.MetricStorer) {
	metrics := []models.Metric{
		{ID: "gauge1", MType: "gauge", Value: float64Pointer(0.75)},
		{ID: "counter1", MType: "counter", Delta: int64Pointer(10)},
	}

	err := stor.StoreSlice(context.Background(), metrics)
	if err != nil {
		t.Fatalf("returned an unexpected error: %v", err)
	}

	for _, metric := range metrics {
		storedMetric, exists, getErr := stor.Get(context.Background(), metric.ID)

		if getErr != nil {
			t.Fatalf("returned an unexpected error: %v", getErr)
		}

		if !exists {
			t.Errorf("expected metric %s not found", metric.ID)
		}

		if !compareMetrics(storedMetric, metric) {
			t.Errorf("Expected metric to be %v, got %v", metric, storedMetric)
		}
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

// Utility functions to create pointers to basic types
func int64Pointer(i int64) *int64 {
	return &i
}

func float64Pointer(f float64) *float64 {
	return &f
}
