package storage

import (
	"fmt"

	"github.com/VOTONO/go-metrics/internal/models"
)

type Storage interface {
	Store(metric models.Metric) error
	Get(name string) (models.Metric, bool)
	All() map[string]models.Metric
}

type StorageImpl struct {
	metrics map[string]models.Metric
}

func New(storage map[string]models.Metric) *StorageImpl {
	if storage == nil {
		storage = make(map[string]models.Metric)
	}
	return &StorageImpl{
		metrics: storage,
	}
}

func (m *StorageImpl) Store(metric models.Metric) error {
	switch metric.Type {
	case "gauge":
		m.metrics[metric.Name] = metric
		fmt.Println("Replaced metric", m.metrics[metric.Name], "with", metric)
		return nil
	case "counter":
		if _, found := m.metrics[metric.Name]; !found {
			m.metrics[metric.Name] = metric
			return nil
		}

		oldMetric := m.metrics[metric.Name]
		newMetric, err := oldMetric.Add(metric)
		if err != nil {
			return fmt.Errorf("error adding metric: %e", err)
		} else {
			m.metrics[metric.Name] = newMetric
			fmt.Printf("Update metric name %s: with value: %v\n", metric.Name, metric.Value)
			return nil
		}

	default:
		return fmt.Errorf("unsupported metric type: %T", metric.Type)
	}
}

func (m *StorageImpl) Get(name string) (models.Metric, bool) {
	metric, found := m.metrics[name]
	return metric, found
}

func (m *StorageImpl) All() map[string]models.Metric {
	return m.metrics
}
