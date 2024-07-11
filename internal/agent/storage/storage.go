package storage

import (
	"sync"

	"github.com/VOTONO/go-metrics/internal/models"
)

type StorageImpl struct {
	mu      sync.Mutex
	metrics map[string]models.Metric
}

func New() *StorageImpl {
	return &StorageImpl{
		metrics: make(map[string]models.Metric),
	}
}

// Returns a copy of the entire map of metrics.
func (s *StorageImpl) Get() map[string]models.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics := make(map[string]models.Metric)
	for metric, value := range s.metrics {
		metrics[metric] = value
	}
	return metrics
}

// Sets the given map into the metrics map.
func (s *StorageImpl) Set(metrics map[string]models.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range metrics {
		s.metrics[key] = value
	}
}
