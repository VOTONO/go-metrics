package repo

import (
	"context"
	"github.com/VOTONO/go-metrics/internal/models"
)

type MetricStorer interface {
	// StoreSingle metric and return updated version
	StoreSingle(ctx context.Context, metric models.Metric) (*models.Metric, error)
	// StoreSlice of metrics
	StoreSlice(ctx context.Context, metrics []models.Metric) error
	// Get metric by ID
	Get(ctx context.Context, ID string) (models.Metric, bool)
	// All stored metrics
	All(ctx context.Context) (map[string]models.Metric, error)
}
