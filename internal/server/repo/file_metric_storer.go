package repo

import (
	"context"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
	"sync"
)

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

func (s *FileMetricStorerImpl) StoreSingle(ctx context.Context, newMetric models.Metric) (*models.Metric, error) {
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

func (s *FileMetricStorerImpl) StoreSlice(ctx context.Context, newMetrics []models.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics, err := ReadFile(s.filePath, s.zapLogger)

	if err != nil {
		return err
	}

	for _, metric := range newMetrics {
		_, err := helpers.UpdateMetricInMap(&metrics, metric, s.zapLogger)
		if err != nil {
			return err
		}
	}

	RewriteFile(s.filePath, metrics, s.zapLogger)

	return nil
}

func (s *FileMetricStorerImpl) Get(ctx context.Context, ID string) (models.Metric, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics, err := ReadFile(s.filePath, s.zapLogger)

	if err != nil {
		return models.Metric{}, false, err
	}

	metric, found := metrics[ID]

	return metric, found, nil
}

func (s *FileMetricStorerImpl) All(ctx context.Context) (map[string]models.Metric, error) {
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
