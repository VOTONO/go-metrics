package storage

import (
	"fmt"
	"sync"

	"github.com/VOTONO/go-metrics/internal/models"
)

type Storage interface {
	Store(metric models.Metric) (*models.Metric, error)
	Get(ID string) (models.Metric, bool)
	All() map[string]models.Metric
}

type StorageImpl struct {
	mu      sync.RWMutex
	metrics map[string]models.Metric
}

func New(initialStorage map[string]models.Metric) *StorageImpl {
	if initialStorage == nil {
		initialStorage = make(map[string]models.Metric)
	}
	return &StorageImpl{
		metrics: initialStorage,
	}
}

func (s *StorageImpl) Store(metric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric.MType {
	case "gauge":
		s.metrics[metric.ID] = metric
		return &metric, nil
	case "counter":
		if existingMetric, found := s.metrics[metric.ID]; found {
			if existingMetric.Delta != nil && metric.Delta != nil {
				*existingMetric.Delta += *metric.Delta
				s.metrics[metric.ID] = existingMetric
				return &existingMetric, nil
			}
			return nil, fmt.Errorf("existing or new metric delta is nil (existing: %v, new: %v)", existingMetric.Delta, metric.Delta)
		}
		s.metrics[metric.ID] = metric
		return &metric, nil
	default:
		return nil, fmt.Errorf("unsupported metric type: %s", metric.MType)
	}
}

func (s *StorageImpl) Get(ID string) (models.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, found := s.metrics[ID]
	return metric, found
}

func (s *StorageImpl) All() map[string]models.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metricsCopy := make(map[string]models.Metric)
	for k, v := range s.metrics {
		metricsCopy[k] = v
	}
	return metricsCopy
}
