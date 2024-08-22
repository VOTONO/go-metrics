package models

import (
	"fmt"
	"strconv"

	"github.com/VOTONO/go-metrics/internal/constants"
)

// Metric struct defines a metric with an ID, type, and either a Delta or Value
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetric is a factory method to create a new Metric based on the type and value provided
func NewMetric(id string, metricType string, value string) (Metric, error) {
	var metric Metric
	var err error

	switch metricType {
	case constants.Gauge:
		var v float64
		v, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid metric value: %v", err)
		}
		metric = Metric{
			ID:    id,
			MType: constants.Gauge,
			Value: &v,
		}
	case constants.Counter:
		var v int64
		v, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid metric delta: %v", err)
		}
		metric = Metric{
			ID:    id,
			MType: constants.Counter,
			Delta: &v,
		}
	default:
		return Metric{}, fmt.Errorf("invalid metric type")
	}

	return metric, nil
}
