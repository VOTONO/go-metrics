// Package repo contains MetricStorer interface and different implementations.
package repo

import (
	"context"

	"github.com/VOTONO/go-metrics/internal/models"
)

// MetricStorer interface for working with metrics storage.
type MetricStorer interface {
	// StoreSingle stores single metric and return updated version.
	StoreSingle(ctx context.Context, metric models.Metric) (*models.Metric, error)
	// StoreSlice sores slice of metrics.
	StoreSlice(ctx context.Context, metrics []models.Metric) error
	// Get return metric by ID, if it exists.
	Get(ctx context.Context, ID string) (models.Metric, bool, error)
	// All returns all stored metrics.
	All(ctx context.Context) (map[string]models.Metric, error)
}
