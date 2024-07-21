package storage

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/VOTONO/go-metrics/internal/models"
)

type MetricStorer interface {
	Store(metric models.Metric) (*models.Metric, error)
	Get(ID string) (models.Metric, bool)
	All() map[string]models.Metric
	Ping() error
}

type MetricStorerImpl struct {
	mu      sync.RWMutex
	db      *sql.DB
	zap     zap.SugaredLogger
	metrics map[string]models.Metric
}

func New(initialStorage map[string]models.Metric, db *sql.DB, zap zap.SugaredLogger) *MetricStorerImpl {
	if initialStorage == nil {
		initialStorage = make(map[string]models.Metric)
	}
	return &MetricStorerImpl{
		metrics: initialStorage,
		db:      db,
		zap:     zap,
	}
}

func (s *MetricStorerImpl) Store(metric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric.MType {
	case "gauge":
		s.metrics[metric.ID] = metric
		s.zap.Infow(
			"new stored metric",
			"metric_id", metric.ID,
			"metric_type", metric.MType,
			"metric_value", metric.Value,
		)
		return &metric, nil
	case "counter":
		if existingMetric, found := s.metrics[metric.ID]; found {
			if existingMetric.Delta != nil && metric.Delta != nil {
				*existingMetric.Delta += *metric.Delta
				s.metrics[metric.ID] = existingMetric
				s.zap.Infow(
					"updated metric",
					"metric_id", existingMetric.ID,
					"metric_type", existingMetric.MType,
					"metric_delta", existingMetric.Delta,
				)
				return &existingMetric, nil
			}
			s.zap.Errorw("existing or new metric delta is nil (existing: %v, new: %v)", existingMetric.Delta, metric.Delta)
			return nil, fmt.Errorf("existing or new metric delta is nil (existing: %v, new: %v)", existingMetric.Delta, metric.Delta)
		}
		s.metrics[metric.ID] = metric
		return &metric, nil
	default:
		s.zap.Errorw("unsupported metric type: %s", metric.MType)
		return nil, fmt.Errorf("unsupported metric type: %s", metric.MType)
	}
}

func (s *MetricStorerImpl) Get(ID string) (models.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, found := s.metrics[ID]
	s.zap.Infow(
		"get metric",
		"found", found,
		"metric_id", ID,
	)
	return metric, found
}

func (s *MetricStorerImpl) All() map[string]models.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metricsCopy := make(map[string]models.Metric)
	for k, v := range s.metrics {
		metricsCopy[k] = v
	}
	s.zap.Infow(
		"all metrics",
		"metrics", metricsCopy,
	)
	return metricsCopy
}

func (s *MetricStorerImpl) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}
