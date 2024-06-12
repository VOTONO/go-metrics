package storage

import (
	"sync"

	"github.com/VOTONO/go-metrics/internal/models"
)

// StorageImpl is the storage structure for metrics.
type StorageImpl struct {
	mu   sync.Mutex
	data map[string]models.Metric
}

// NewStorage creates and returns a new StorageImpl instance.
func New(metrics map[string]models.Metric) *StorageImpl {
	return &StorageImpl{
		data: make(map[string]models.Metric),
	}
}

// Get returns a copy of the entire map of metrics.
func (s *StorageImpl) Get() map[string]models.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := make(map[string]models.Metric)
	for key, value := range s.data {
		copy[key] = value
	}
	return copy
}

// Set sets the values from the given map in the metrics map.
func (s *StorageImpl) Set(metrics map[string]models.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range metrics {
		s.data[key] = value
	}
}
