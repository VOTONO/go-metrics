package utils

import "github.com/VOTONO/go-metrics/internal/models"

var ValidFloat64 = float64(456.78)
var ValidInt64 = int64(456)

var ValidGaugeMetric = models.Metric{
	ID:    "validGaugeMetric",
	MType: "gauge",
	Value: &ValidFloat64,
}

var ValidCounterMetric = models.Metric{
	ID:    "validCounterMetric",
	MType: "counter",
	Delta: &ValidInt64,
}

// Invalid Metrics
var InvalidGaugeMissingValue = models.Metric{
	ID:    "invalidGaugeMissingValue",
	MType: "gauge",
	Value: nil, // Missing value for gauge
}

var InvalidCounterMissingDelta = models.Metric{
	ID:    "invalidCounterMissingDelta",
	MType: "counter",
	Delta: nil, // Missing delta for counter
}

var InvalidMetricUnknownType = models.Metric{
	ID:    "invalidMetricUnknownType",
	MType: "unknown", // Unknown type
	Value: &ValidFloat64,
}

var InvalidMetricCounterWithGaugeValue = models.Metric{
	ID:    "invalidMetricCounterWithGaugeValue",
	MType: "counter",
	Value: &ValidFloat64, // Counter metric should not have a gauge value
}

var InvalidMetricGaugeWithCounterDelta = models.Metric{
	ID:    "invalidMetricGaugeWithCounterDelta",
	MType: "gauge",
	Delta: &ValidInt64, // Gauge metric should not have a counter delta
}

var InvalidMetricEmptyID = models.Metric{
	ID:    "",
	MType: "gauge",
	Value: &ValidFloat64, // Empty ID
}

var InvalidMetricNilBothFields = models.Metric{
	ID:    "invalidMetricNilBothFields",
	MType: "gauge",
	Value: nil,
	Delta: nil, // Both fields are nil
}
