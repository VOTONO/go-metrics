package repo

import (
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
	"sync"
)

type MetricStorer interface {
	Store(metric models.Metric) (*models.Metric, error)
	Get(ID string) (models.Metric, bool)
	All() (map[string]models.Metric, error)
	Ping() error
}

type LocalMetricStorerImpl struct {
	mu        sync.RWMutex
	zapLogger *zap.SugaredLogger
	metrics   map[string]models.Metric
}

func NewLocalMetricStorer(restore bool, filePath string, zapLogger *zap.SugaredLogger) *LocalMetricStorerImpl {
	storer := &LocalMetricStorerImpl{
		metrics:   make(map[string]models.Metric),
		zapLogger: zapLogger,
	}
	if restore {
		restoredMetrics, err := ReadFile(filePath, zapLogger)
		if err != nil {
			return storer
		}
		storer.metrics = restoredMetrics
	}
	return storer
}

func (s *LocalMetricStorerImpl) Store(newMetric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updatedMetric, err := helpers.UpdateMetricInMap(&s.metrics, newMetric, s.zapLogger)

	if err != nil {
		return &models.Metric{}, err
	}

	helpers.LogMetric("new stored Metric", updatedMetric, s.zapLogger)
	return &updatedMetric, nil
}

func (s *LocalMetricStorerImpl) Get(ID string) (models.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, found := s.metrics[ID]

	return metric, found
}

func (s *LocalMetricStorerImpl) All() (map[string]models.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metricsCopy := make(map[string]models.Metric)
	for k, v := range s.metrics {
		metricsCopy[k] = v
	}

	return metricsCopy, nil
}

func (s *LocalMetricStorerImpl) Ping() error {
	return nil
}
