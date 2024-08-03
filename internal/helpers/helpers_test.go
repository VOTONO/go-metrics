package helpers

import (
	"github.com/VOTONO/go-metrics/internal/models"
	"testing"
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
			args:    models.Metric{ID: "metric1", MType: "gauge", Value: float64Ptr(123.45)},
			want:    "123.45",
			wantErr: false,
		},
		{
			name:    "Extract counter value",
			args:    models.Metric{ID: "metric2", MType: "counter", Delta: int64Ptr(123)},
			want:    "123",
			wantErr: false,
		},
		{
			name:    "Gauge value is nil",
			args:    models.Metric{ID: "metric3", MType: "gauge", Value: nil},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Counter delta is nil",
			args:    models.Metric{ID: "metric4", MType: "counter", Delta: nil},
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
			args: models.Metric{ID: "metric1", MType: "gauge", Value: float64Ptr(123.45)},
			want: true,
		},
		{
			name: "Valid counter metric",
			args: models.Metric{ID: "metric2", MType: "counter", Delta: int64Ptr(123)},
			want: true,
		},
		{
			name: "Gauge metric with nil value",
			args: models.Metric{ID: "metric3", MType: "gauge", Value: nil},
			want: false,
		},
		{
			name: "Counter metric with nil delta",
			args: models.Metric{ID: "metric4", MType: "counter", Delta: nil},
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
			old:     models.Metric{ID: "metric1", MType: "counter", Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric1", MType: "counter", Delta: int64Ptr(5)},
			want:    models.Metric{ID: "metric1", MType: "counter", Delta: int64Ptr(15)},
			wantErr: false,
		},
		{
			name:    "Metric type mismatch",
			old:     models.Metric{ID: "metric2", MType: "counter", Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric2", MType: "gauge", Value: float64Ptr(123.45)},
			want:    models.Metric{ID: "metric2", MType: "counter", Delta: int64Ptr(10)},
			wantErr: true,
		},
		{
			name:    "Old delta is nil",
			old:     models.Metric{ID: "metric3", MType: "counter", Delta: nil},
			new:     models.Metric{ID: "metric3", MType: "counter", Delta: int64Ptr(5)},
			want:    models.Metric{ID: "metric3", MType: "counter", Delta: nil},
			wantErr: true,
		},
		{
			name:    "New delta is nil",
			old:     models.Metric{ID: "metric4", MType: "counter", Delta: int64Ptr(10)},
			new:     models.Metric{ID: "metric4", MType: "counter", Delta: nil},
			want:    models.Metric{ID: "metric4", MType: "counter", Delta: int64Ptr(10)},
			wantErr: true,
		},
		{
			name:    "Both deltas are nil",
			old:     models.Metric{ID: "metric5", MType: "counter", Delta: nil},
			new:     models.Metric{ID: "metric5", MType: "counter", Delta: nil},
			want:    models.Metric{ID: "metric5", MType: "counter", Delta: nil},
			wantErr: true,
		},
		{
			name:    "Different names",
			old:     models.Metric{ID: "old", MType: "counter", Delta: nil},
			new:     models.Metric{ID: "new", MType: "counter", Delta: nil},
			want:    models.Metric{ID: "metric5", MType: "counter", Delta: nil},
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
