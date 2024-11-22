package helpers

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

func float64Ptr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64       { return &v }

func TestExtractValue(t *testing.T) {
	tests := []struct {
		name    string
		args    models.Metric
		want    string
		wantErr bool
	}{
		{
			name:    "Extract gauge value",
			args:    models.Metric{ID: "metric1", MType: constants.Gauge, Value: float64Ptr(123.45)},
			want:    "123.45",
			wantErr: false,
		},
		{
			name:    "Extract counter value",
			args:    models.Metric{ID: "metric2", MType: constants.Counter, Delta: int64Ptr(123)},
			want:    "123",
			wantErr: false,
		},
		{
			name:    "Gauge value is nil",
			args:    models.Metric{ID: "metric3", MType: constants.Gauge, Value: nil},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Counter delta is nil",
			args:    models.Metric{ID: "metric4", MType: constants.Counter, Delta: nil},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Unknown metric type",
			args:    models.Metric{ID: "metric5", MType: "unknown"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractValue(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMetric(t *testing.T) {
	tests := []struct {
		name string
		args models.Metric
		want bool
	}{
		{
			name: "Valid gauge metric",
			args: models.Metric{ID: "metric1", MType: constants.Gauge, Value: float64Ptr(123.45)},
			want: true,
		},
		{
			name: "Valid counter metric",
			args: models.Metric{ID: "metric2", MType: constants.Counter, Delta: int64Ptr(123)},
			want: true,
		},
		{
			name: "Gauge metric with nil value",
			args: models.Metric{ID: "metric3", MType: constants.Gauge, Value: nil},
			want: false,
		},
		{
			name: "Counter metric with nil delta",
			args: models.Metric{ID: "metric4", MType: constants.Counter, Delta: nil},
			want: false,
		},
		{
			name: "Unknown metric type",
			args: models.Metric{ID: "metric5", MType: "unknown"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateMetric(tt.args); got != tt.want {
				t.Errorf("ValidateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	tests := []struct {
		name    string
		old     models.Metric
		new     models.Metric
		want    models.Metric
		wantErr bool
	}{
		{
			name:    "Valid counter metric update",
			old:     models.Metric{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(5)},
			want:    models.Metric{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(15)},
			wantErr: false,
		},
		{
			name:    "Metric type mismatch",
			old:     models.Metric{ID: "metric2", MType: constants.Counter, Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(123.45)},
			want:    models.Metric{ID: "metric2", MType: constants.Counter, Delta: int64Ptr(10)},
			wantErr: true,
		},
		{
			name:    "New delta is nil",
			old:     models.Metric{ID: "metric4", MType: constants.Counter, Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric4", MType: constants.Counter, Delta: nil},
			want:    models.Metric{ID: "metric4", MType: constants.Counter, Delta: int64Ptr(10)},
			wantErr: true,
		},
		{
			name:    "Both deltas are nil",
			old:     models.Metric{ID: "metric5", MType: constants.Counter, Delta: nil},
			new:     models.Metric{ID: "metric5", MType: constants.Counter, Delta: nil},
			want:    models.Metric{ID: "metric5", MType: constants.Counter, Delta: nil},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateCounterMetric(tt.old, tt.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCounterMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID || got.MType != tt.want.MType || (got.Delta == nil && tt.want.Delta != nil) || (got.Delta != nil && tt.want.Delta == nil) || (got.Delta != nil && tt.want.Delta != nil && *got.Delta != *tt.want.Delta) || (got.Value == nil && tt.want.Value != nil) || (got.Value != nil && tt.want.Value == nil) || (got.Value != nil && tt.want.Value != nil && *got.Value != *tt.want.Value) {
				t.Errorf("UpdateCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsToHTML(t *testing.T) {
	tests := []struct {
		name      string
		metrics   map[string]models.Metric
		want      string
		expectErr bool
	}{
		{
			name: "Valid metrics",
			metrics: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: constants.Gauge, Value: func() *float64 { f := 123.45; return &f }()},
				"metric2": {ID: "metric2", MType: constants.Counter, Delta: func() *int64 { i := int64(678); return &i }()},
			},
			want:      `<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr><tr><td>metric1</td><td>123.45</td></tr><tr><td>metric2</td><td>678</td></tr></table></body></html>`,
			expectErr: false,
		},
		{
			name: "constants.Gauge metric with nil value",
			metrics: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: constants.Gauge, Value: nil},
			},
			want:      "",
			expectErr: true,
		},
		{
			name: "Counter metric with nil delta",
			metrics: map[string]models.Metric{
				"metric2": {ID: "metric2", MType: constants.Counter, Delta: nil},
			},
			want:      "",
			expectErr: true,
		},
		{
			name: "Unsupported metric type",
			metrics: map[string]models.Metric{
				"metric3": {ID: "metric3", MType: "unknown"},
			},
			want:      "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}
			defer logger.Sync()
			zapLogger := *logger.Sugar()
			got, err := MetricsToHTML(tt.metrics, &zapLogger)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func BenchmarkExtractValue(b *testing.B) {
	metricGauge := models.Metric{
		ID:    "gaugeMetric",
		MType: constants.Gauge,
		Value: float64Ptr(123.456),
	}

	metricCounter := models.Metric{
		ID:    "counterMetric",
		MType: constants.Counter,
		Delta: int64Ptr(789),
	}

	for i := 0; i < b.N; i++ {
		_, _ = ExtractValue(metricGauge)
		_, _ = ExtractValue(metricCounter)
	}
}

func BenchmarkValidateMetric(b *testing.B) {
	metricGauge := models.Metric{
		ID:    "gaugeMetric",
		MType: constants.Gauge,
		Value: float64Ptr(123.456),
	}

	metricCounter := models.Metric{
		ID:    "counterMetric",
		MType: constants.Counter,
		Delta: int64Ptr(789),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateMetric(metricGauge)
		_ = ValidateMetric(metricCounter)
	}
}

func BenchmarkUpdateCounterMetric(b *testing.B) {
	oldMetric := models.Metric{
		ID:    "counterMetric",
		MType: constants.Counter,
		Delta: int64Ptr(100),
	}

	newMetric := models.Metric{
		ID:    "counterMetric",
		MType: constants.Counter,
		Delta: int64Ptr(50),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = UpdateCounterMetric(oldMetric, newMetric)
	}
}

func BenchmarkUpdateMetricInMap(b *testing.B) {
	logger := zap.NewNop().Sugar()

	metrics := map[string]models.Metric{
		"gaugeMetric": {
			ID:    "gaugeMetric",
			MType: constants.Gauge,
			Value: float64Ptr(123.456),
		},
		"counterMetric": {
			ID:    "counterMetric",
			MType: constants.Counter,
			Delta: int64Ptr(100),
		},
	}

	newMetric := models.Metric{
		ID:    "counterMetric",
		MType: constants.Counter,
		Delta: int64Ptr(50),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = UpdateMetricInMap(metrics, newMetric, logger)
	}
}

func BenchmarkProcessMetricsDuplicates(b *testing.B) {
	metrics := []models.Metric{
		{ID: "counterMetric", MType: constants.Counter, Delta: int64Ptr(100)},
		{ID: "counterMetric", MType: constants.Counter, Delta: int64Ptr(50)},
		{ID: "gaugeMetric", MType: constants.Gauge, Value: float64Ptr(123.456)},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ProcessMetricsDuplicates(metrics)
	}
}

func TestConvertMapToSlice(t *testing.T) {
	type args struct {
		metricsMap map[string]models.Metric
	}
	tests := []struct {
		name string
		args args
		want []models.Metric
	}{
		{
			name: "Single metric",
			args: args{
				metricsMap: map[string]models.Metric{
					"metric1_counter": {ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
				},
			},
			want: []models.Metric{{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)}},
		},
		{
			name: "Multiple metrics",
			args: args{
				metricsMap: map[string]models.Metric{
					"metric1_counter": {ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
					"metric2_gauge":   {ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)},
				},
			},
			want: []models.Metric{
				{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
				{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, ConvertMapToSlice(tt.args.metricsMap), "ConvertMapToSlice(%v)", tt.args.metricsMap)
		})
	}
}

func TestProcessMetricsDuplicates(t *testing.T) {
	type args struct {
		metrics []models.Metric
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "No duplicates",
			args: args{
				metrics: []models.Metric{
					{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
					{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)},
				},
			},
			want: []models.Metric{
				{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
				{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Duplicate counter metrics",
			args: args{
				metrics: []models.Metric{
					{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
					{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(15)},
				},
			},
			want: []models.Metric{
				{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(25)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Mixed duplicates",
			args: args{
				metrics: []models.Metric{
					{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(10)},
					{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(5)},
					{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(2.71)},
					{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)}, // Gauge should take the last seen
				},
			},
			want: []models.Metric{
				{ID: "metric1", MType: constants.Counter, Delta: int64Ptr(15)},
				{ID: "metric2", MType: constants.Gauge, Value: float64Ptr(3.14)},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessMetricsDuplicates(tt.args.metrics)
			if !tt.wantErr(t, err, fmt.Sprintf("ProcessMetricsDuplicates(%v)", tt.args.metrics)) {
				return
			}
			assert.ElementsMatch(t, tt.want, got, "ProcessMetricsDuplicates(%v)", tt.args.metrics)
		})
	}
}
