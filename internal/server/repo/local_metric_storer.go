package repo

import (
	"context"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
	"sync"
)

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

func (s *LocalMetricStorerImpl) StoreSingle(ctx context.Context, newMetric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valid := helpers.ValidateMetric(newMetric)
	if !valid {
		return &models.Metric{}, fmt.Errorf("invalide metric")
	}

	updatedMetric, err := helpers.UpdateMetricInMap(&s.metrics, newMetric, s.zapLogger)

	if err != nil {
		return &models.Metric{}, err
	}

	helpers.LogMetric("new stored Metric", updatedMetric, s.zapLogger)
	return &updatedMetric, nil
}

func (s *LocalMetricStorerImpl) StoreSlice(ctx context.Context, newMetrics []models.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, metric := range newMetrics {
		_, err := helpers.UpdateMetricInMap(&s.metrics, metric, s.zapLogger)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *LocalMetricStorerImpl) Get(ctx context.Context, ID string) (models.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metric, found := s.metrics[ID]

	return metric, found
}

func (s *LocalMetricStorerImpl) All(ctx context.Context) (map[string]models.Metric, error) {
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
