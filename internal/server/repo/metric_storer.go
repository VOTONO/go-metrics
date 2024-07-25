package repo

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
	mu            sync.RWMutex
	db            *sql.DB
	zapLogger     *zap.SugaredLogger
	metrics       map[string]models.Metric
	filePath      string
	storeInterval int
}

func New(restore bool, filePath string, storeInterval int, db *sql.DB, zapLogger *zap.SugaredLogger) *MetricStorerImpl {
	storer := &MetricStorerImpl{
		metrics:       make(map[string]models.Metric),
		filePath:      filePath,
		storeInterval: storeInterval,
		db:            db,
		zapLogger:     zapLogger,
	}
	if restore {
		restoredMetrics, err := Read(filePath, zapLogger)
		if err != nil {
			zapLogger.Errorw("Fail read metrics from file", "path", filePath, "error", err)
			return storer
		}
		storer.metrics = restoredMetrics
	}
	return storer
}

func (s *MetricStorerImpl) Store(metric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metric.MType {
	case "gauge":
		s.metrics[metric.ID] = metric
		s.zapLogger.Infow(
			"new stored gauge metric",
			"metric_id", metric.ID,
			"metric_type", metric.MType,
			"metric_value", metric.Value,
		)
	case "counter":
		if existingMetric, found := s.metrics[metric.ID]; found {
			if existingMetric.Delta != nil {
				*existingMetric.Delta += *metric.Delta
				s.metrics[metric.ID] = existingMetric
				s.zapLogger.Infow(
					"updated counter metric",
					"metric_id", existingMetric.ID,
					"metric_type", existingMetric.MType,
					"metric_delta", existingMetric.Delta,
				)
			} else {
				err := fmt.Errorf("existing metric delta is nil (existing: %v, new: %v)", existingMetric.Delta, metric.Delta)
				s.zapLogger.Errorw(
					"error updating counter metric",
					"metric_id", metric.ID,
					"existing_delta", existingMetric.Delta,
					"new_delta", metric.Delta,
					"error", err,
				)
				return nil, err
			}
		} else {
			s.metrics[metric.ID] = metric
			s.zapLogger.Infow(
				"new stored counter metric",
				"metric_id", metric.ID,
				"metric_type", metric.MType,
				"metric_delta", metric.Delta,
			)
		}
	default:
		err := fmt.Errorf("unsupported metric type: %s", metric.MType)
		s.zapLogger.Errorw("error storing metric", "metric_id", metric.ID, "error", err)
		return nil, err
	}

	if s.storeInterval == 0 && s.filePath != "" {
		if err := Write(s.filePath, s.metrics, s.zapLogger); err != nil {
			s.zapLogger.Errorw("error writing metrics to file", "file_path", s.filePath, "error", err)
			return nil, err
		}
	}

	m := s.metrics[metric.ID]
	return &m, nil
}

func (s *MetricStorerImpl) Get(ID string) (models.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, found := s.metrics[ID]
	s.zapLogger.Infow(
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
	s.zapLogger.Infow(
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
