package repo

import (
	"sync"

	"github.com/VOTONO/go-metrics/internal/models"
)

// MetricStorer defines the interface for a metric storage.
type MetricStorer interface {
	Get() map[string]models.Metric
	Set(metrics map[string]models.Metric)
}

// MetricStorerImpl is a thread-safe implementation of MetricStorer.
type MetricStorerImpl struct {
	mu      sync.RWMutex
	metrics map[string]models.Metric
}

// New creates a new instance of MetricStorerImpl.
func New() MetricStorer {
	return &MetricStorerImpl{
		metrics: make(map[string]models.Metric),
	}
}

// Get returns a copy of the entire map of metrics.
func (s *MetricStorerImpl) Get() map[string]models.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metricsCopy := make(map[string]models.Metric, len(s.metrics))
	for key, value := range s.metrics {
		metricsCopy[key] = value
	}
	return metricsCopy
}

// Set stores the given map into the metrics map.
func (s *MetricStorerImpl) Set(metrics map[string]models.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range metrics {
		s.metrics[key] = value
	}
}
