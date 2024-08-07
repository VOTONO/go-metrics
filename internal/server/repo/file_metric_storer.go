package repo

import (
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

func (s *FileMetricStorerImpl) Store(newMetric models.Metric) (*models.Metric, error) {
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

func (s *FileMetricStorerImpl) Get(ID string) (models.Metric, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics, err := ReadFile(s.filePath, s.zapLogger)

	if err != nil {
		return models.Metric{}, false
	}

	metric, found := metrics[ID]

	return metric, found
}

func (s *FileMetricStorerImpl) All() (map[string]models.Metric, error) {
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
