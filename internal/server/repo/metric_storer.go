package repo

import "github.com/VOTONO/go-metrics/internal/models"

type MetricStorer interface {
	// Store metric and return updated version
	Store(metric models.Metric) (*models.Metric, error)
	// Get metric by ID
	Get(ID string) (models.Metric, bool)
	// All stored metrics
	All() (map[string]models.Metric, error)
}
