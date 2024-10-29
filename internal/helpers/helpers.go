package helpers

import (
	"errors"
	"fmt"
	"html"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

// ExtractValue converts metric value to a string based on its type
func ExtractValue(m models.Metric) (string, error) {
	switch m.MType {
	case constants.Gauge:
		if m.Value != nil {
			return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
		}
		return "", errors.New("metric value not found")
	case constants.Counter:
		if m.Delta != nil {
			return strconv.FormatInt(*m.Delta, 10), nil
		}
		return "", errors.New("metric delta not found")
	default:
		return "", errors.New("unknown metric type")
	}
}

// ValidateMetric checks if a metric is valid based on its type and fields
func ValidateMetric(m models.Metric) bool {
	if m.ID == "" {
		return false
	}
	return (m.MType == constants.Gauge && m.Value != nil) ||
		(m.MType == constants.Counter && m.Delta != nil)
}

// UpdateCounterMetric aggregates Counter metrics
func UpdateCounterMetric(old, new models.Metric) (models.Metric, error) {
	if old.ID != new.ID || old.MType != constants.Counter || new.MType != constants.Counter {
		return old, errors.New("metric mismatch")
	}
	if old.Delta == nil || new.Delta == nil {
		return old, errors.New("metric delta missing")
	}
	newDelta := *old.Delta + *new.Delta
	new.Delta = &newDelta
	return new, nil
}

// UpdateMetricInMap updates a metric in the given map, creating or updating as needed
func UpdateMetricInMap(metrics map[string]models.Metric, metric models.Metric, logger *zap.SugaredLogger) (models.Metric, error) {
	if metric.MType == constants.Gauge {
		metrics[metric.ID] = metric
		return metric, nil
	}

	if metric.MType == constants.Counter {
		if existingMetric, found := metrics[metric.ID]; found {
			updatedMetric, err := UpdateCounterMetric(existingMetric, metric)
			if err != nil {
				logger.Errorw("fail to update counter Metric", "metric_id", metric.ID, "error", err.Error())
				return models.Metric{}, err
			}
			metrics[metric.ID] = updatedMetric
			return updatedMetric, nil
		}
		metrics[metric.ID] = metric
		return metric, nil
	}

	err := fmt.Errorf("unsupported Metric type: %s", metric.MType)
	logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
	return models.Metric{}, err
}

// MetricsToHTML generates an HTML table of metrics, using a builder for performance
func MetricsToHTML(metrics map[string]models.Metric, logger *zap.SugaredLogger) (string, error) {
	keys := make([]string, 0, len(metrics))
	for key := range metrics {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>")

	for _, key := range keys {
		metric := metrics[key]
		value, err := ExtractValue(metric)
		if err != nil {
			logger.Errorw("Invalid metric value", "metric_id", key, "error", err)
			return "", fmt.Errorf("invalid metric value for metric %s: %w", key, err)
		}
		htmlBuilder.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", html.EscapeString(key), html.EscapeString(value)))
	}

	htmlBuilder.WriteString("</table></body></html>")
	return htmlBuilder.String(), nil
}

// ProcessMetricsDuplicates consolidates duplicate metrics in a slice by ID and type
func ProcessMetricsDuplicates(metrics []models.Metric) ([]models.Metric, error) {
	metricMap := make(map[string]models.Metric)

	for _, metric := range metrics {
		key := metric.ID + "_" + metric.MType

		if existingMetric, exists := metricMap[key]; exists && metric.MType == constants.Counter {
			updatedMetric, err := UpdateCounterMetric(existingMetric, metric)
			if err != nil {
				return metrics, err
			}
			metricMap[key] = updatedMetric
		} else {
			metricMap[key] = metric
		}
	}

	return ConvertMapToSlice(metricMap), nil
}

// ConvertMapToSlice converts a map to a slice of metrics
func ConvertMapToSlice(metricsMap map[string]models.Metric) []models.Metric {
	metricsSlice := make([]models.Metric, 0, len(metricsMap))
	for _, metric := range metricsMap {
		metricsSlice = append(metricsSlice, metric)
	}
	return metricsSlice
}

// LogMetric logs a metric's details
func LogMetric(message string, metric models.Metric, logger *zap.SugaredLogger) {
	logger.Infow(
		message,
		"metric_id", metric.ID,
		"metric_type", metric.MType,
		"metric_value", metric.Value,
		"metric_delta", metric.Delta,
	)
}

// DecideShouldRetryAfterError checks if an error is retryable (database or file-related)
func DecideShouldRetryAfterError(err error) bool {
	var pgErr *pgconn.PgError
	var pathErr *os.PathError

	return (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException) ||
		(errors.As(err, &pathErr) && errors.Is(pathErr.Err, syscall.EBUSY))
}
