package repo

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
)

// FileMetricStorerImpl implementation of MetricStorer interface. Stores all metrics in file.
type FileMetricStorerImpl struct {
	mu        sync.RWMutex
	filePath  string
	zapLogger *zap.SugaredLogger
}

func NewFileMetricStorer(filePath string, logger *zap.SugaredLogger) MetricStorer {
	return &FileMetricStorerImpl{
		filePath:  filePath,
		zapLogger: logger,
	}
}

func (s *FileMetricStorerImpl) StoreSingle(_ context.Context, newMetric models.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valid := helpers.ValidateMetric(newMetric)
	if !valid {
		return &models.Metric{}, fmt.Errorf("invalide metric")
	}

	updatedMetric, err := AddToFile(s.filePath, newMetric, s.zapLogger)

	if err != nil {
		return nil, err
	}

	helpers.LogMetric("new stored Metric", updatedMetric, s.zapLogger)
	return &updatedMetric, nil
}

func (s *FileMetricStorerImpl) StoreSlice(_ context.Context, newMetrics []models.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics, readErr := ReadFile(s.filePath, s.zapLogger)

	if readErr != nil {
		return readErr
	}

	for _, metric := range newMetrics {
		_, err := helpers.UpdateMetricInMap(metrics, metric, s.zapLogger)
		if err != nil {
			return err
		}
	}

	rewriteErr := RewriteFile(s.filePath, metrics, s.zapLogger)
	if rewriteErr != nil {
		return rewriteErr
	}

	return nil
}

func (s *FileMetricStorerImpl) Get(_ context.Context, id string) (models.Metric, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics, err := ReadFile(s.filePath, s.zapLogger)

	if err != nil {
		return models.Metric{}, false, err
	}

	metric, found := metrics[id]

	return metric, found, nil
}

func (s *FileMetricStorerImpl) All(_ context.Context) (map[string]models.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics, err := ReadFile(s.filePath, s.zapLogger)

	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *FileMetricStorerImpl) Ping() error {
	return nil
}
