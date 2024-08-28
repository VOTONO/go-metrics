package utils

import (
	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

var ValidFloat64 = float64(456.78)
var ValidInt64 = int64(456)

var ValidGaugeMetric = models.Metric{
	ID:    "validGaugeMetric",
	MType: constants.Gauge,
	Value: &ValidFloat64,
}

var ValidCounterMetric = models.Metric{
	ID:    "validCounterMetric",
	MType: constants.Counter,
	Delta: &ValidInt64,
}

var InvalidGaugeMissingValue = models.Metric{
	ID:    "invalidGaugeMissingValue",
	MType: constants.Gauge,
	Value: nil, // Missing value for gauge
}

var InvalidCounterMissingDelta = models.Metric{
	ID:    "invalidCounterMissingDelta",
	MType: constants.Counter,
	Delta: nil, // Missing delta for counter
}

var InvalidMetricUnknownType = models.Metric{
	ID:    "invalidMetricUnknownType",
	MType: "unknown", // Unknown type
	Value: &ValidFloat64,
}

var InvalidMetricCounterWithGaugeValue = models.Metric{
	ID:    "invalidMetricCounterWithGaugeValue",
	MType: constants.Counter,
	Value: &ValidFloat64, // Counter metric should not have a gauge value
}

var InvalidMetricGaugeWithCounterDelta = models.Metric{
	ID:    "invalidMetricGaugeWithCounterDelta",
	MType: constants.Gauge,
	Delta: &ValidInt64, // Gauge metric should not have a counter delta
}

var InvalidMetricEmptyID = models.Metric{
	ID:    "",
	MType: constants.Gauge,
	Value: &ValidFloat64, // Empty ID
}

var InvalidMetricNilBothFields = models.Metric{
	ID:    "invalidMetricNilBothFields",
	MType: constants.Gauge,
	Value: nil,
	Delta: nil, // Both fields are nil
}
