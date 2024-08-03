package helpers

import (
	"fmt"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
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

func UpdateCounterMetric(old models.Metric, new models.Metric) (models.Metric, error) {

	if old.ID != new.ID {
		return old, fmt.Errorf("metrics have different names")
	}
	if old.MType != "counter" || new.MType != "counter" {
		return old, fmt.Errorf("metric type mismatch")
	}

	if old.Delta == nil || new.Delta == nil {
		return old, fmt.Errorf("metric delta mismatch")
	}

	newDelta := *old.Delta + *new.Delta
	new.Delta = &newDelta

	return new, nil
}

// UpdateMetricInMap update given metrics map with given metric
func UpdateMetricInMap(metrics *map[string]models.Metric, metric models.Metric, logger *zap.SugaredLogger) (models.Metric, error) {
	switch metric.MType {
	case "gauge":
		(*metrics)[metric.ID] = metric
		return metric, nil
	case "counter":
		if existingMetric, found := (*metrics)[metric.ID]; found {
			updatedMetric, err := UpdateCounterMetric(existingMetric, metric)

			if err != nil {
				logger.Errorw("fail to update counter Metric", "metric_id", metric.ID, "error", err.Error())
				return models.Metric{}, err
			}

			(*metrics)[metric.ID] = updatedMetric
			return updatedMetric, nil
		} else {
			(*metrics)[metric.ID] = metric
			return metric, nil
		}
	default:
		err := fmt.Errorf("unsupported Metric type: %s", metric.MType)
		logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
		return models.Metric{}, err
	}
}

func LogMetric(message string, metric models.Metric, logger *zap.SugaredLogger) {
	logger.Infow(
		message,
		"metric_id", metric.ID,
		"metric_type", metric.MType,
		"metric_value", metric.Value,
		"metric_delta", metric.Delta,
	)
}
