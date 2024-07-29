package helpers

import (
	"fmt"
	"github.com/VOTONO/go-metrics/internal/models"
	"strconv"
)

func ExtractValue(m models.Metric) (string, error) {
	var value string

	switch m.MType {
	case "gauge":
		if m.Value != nil {
			value = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else {
			return "", fmt.Errorf("metric value not found")
		}
	case "counter":
		if m.Delta != nil {
			value = strconv.FormatInt(*m.Delta, 10)
		} else {
			return "", fmt.Errorf("metric delta not found")
		}
	default:
		return "", fmt.Errorf("unknown metric type")
	}
	return value, nil
}

func ValidateMetric(m models.Metric) bool {
	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return false
		}
		return true
	case "counter":
		if m.Delta == nil {
			return false
		}
		return true
	default:
		return false
	}
}
