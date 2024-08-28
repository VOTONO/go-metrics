package helpers

import (
	"errors"
	"fmt"
	"html"
	"os"
	"sort"
	"strconv"
	"syscall"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

func ExtractValue(m models.Metric) (string, error) {
	var value string

	switch m.MType {
	case constants.Gauge:
		if m.Value != nil {
			value = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else {
			return "", fmt.Errorf("metric value not found")
		}
	case constants.Counter:
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
	if m.ID == "" {
		return false
	}
	switch m.MType {
	case constants.Gauge:
		if m.Value == nil {
			return false
		}
		return true
	case constants.Counter:
		if m.Delta == nil {
			return false
		}
		return true
	default:
		return false
	}
}

func UpdateCounterMetric(old models.Metric, new models.Metric) (models.Metric, error) {
	oldValid := ValidateMetric(old)
	newValid := ValidateMetric(new)

	if !oldValid && newValid {
		return new, nil
	}

	if old.ID != new.ID {
		return old, fmt.Errorf("metrics have different names")
	}
	if old.MType != constants.Counter || new.MType != constants.Counter {
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
func UpdateMetricInMap(metrics map[string]models.Metric, metric models.Metric, logger *zap.SugaredLogger) (models.Metric, error) {
	switch metric.MType {
	case constants.Gauge:
		(metrics)[metric.ID] = metric
		return metric, nil
	case constants.Counter:
		if existingMetric, found := (metrics)[metric.ID]; found {
			updatedMetric, err := UpdateCounterMetric(existingMetric, metric)

			if err != nil {
				logger.Errorw("fail to update counter Metric", "metric_id", metric.ID, "error", err.Error())
				return models.Metric{}, err
			}

			(metrics)[metric.ID] = updatedMetric
			return updatedMetric, nil
		} else {
			(metrics)[metric.ID] = metric
			return metric, nil
		}
	default:
		err := fmt.Errorf("unsupported Metric type: %s", metric.MType)
		logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
		return models.Metric{}, err
	}
}

func MetricsToHTML(metrics map[string]models.Metric, logger *zap.SugaredLogger) (string, error) {
	keys := make([]string, 0, len(metrics))
	for key := range metrics {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	htmlString := "<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>"
	for _, key := range keys {
		metric := metrics[key]
		value, err := ExtractValue(metric)
		if err != nil {
			logger.Errorw("Invalid metric value", "metric_id", key, "error", err)
			return "", fmt.Errorf("invalid metric value for metric %s: %w", key, err)
		}

		htmlString += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", html.EscapeString(key), html.EscapeString(fmt.Sprintf("%v", value)))
	}

	htmlString += "</table></body></html>"
	return htmlString, nil
}

// ProcessMetricsDuplicates with the same ID and type
func ProcessMetricsDuplicates(metrics []models.Metric) ([]models.Metric, error) {
	metricMap := make(map[string]models.Metric)

	for _, metric := range metrics {
		key := metric.ID + "_" + metric.MType

		if existingMetric, exists := metricMap[key]; exists {
			switch metric.MType {
			case constants.Gauge:
				metricMap[key] = metric
			case constants.Counter:
				updatedMetric, err := UpdateCounterMetric(existingMetric, metric)
				if err != nil {
					return metrics, err
				}
				metricMap[key] = updatedMetric
			default:
				err := fmt.Errorf("unsupported Metric type: %s", metric.MType)
				return nil, err
			}

		} else {
			metricMap[key] = metric
		}
	}

	slice := ConvertMapToSlice(metricMap)
	return slice, nil
}

// ConvertMapToSlice converts a map[string]models.Metric to a slice of models.Metric.
func ConvertMapToSlice(metricsMap map[string]models.Metric) []models.Metric {
	metricsSlice := make([]models.Metric, 0, len(metricsMap))
	for _, metric := range metricsMap {
		metricsSlice = append(metricsSlice, metric)
	}
	return metricsSlice
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

func DecideShouldRetryAfterError(err error) bool {
	var pgErr *pgconn.PgError
	var pathErr *os.PathError

	isDBConnectionError := errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException
	isFileBusyError := errors.As(err, &pathErr) && errors.Is(pathErr.Err, syscall.EBUSY)

	return isDBConnectionError || isFileBusyError
}
