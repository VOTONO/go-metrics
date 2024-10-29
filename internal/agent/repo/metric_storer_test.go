package repo

import (
	"reflect"
	"testing"

	"github.com/VOTONO/go-metrics/internal/models"
)

func TestMetricStorerImpl_Get(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]models.Metric
		want   map[string]models.Metric
	}{
		{
			name:   "Get from empty store",
			fields: map[string]models.Metric{},
			want:   map[string]models.Metric{},
		},
		{
			name: "Get from non-empty store",
			fields: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
				"metric2": {ID: "metric2", MType: "counter", Delta: new(int64)},
			},
			want: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
				"metric2": {ID: "metric2", MType: "counter", Delta: new(int64)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorerImpl{
				metrics: tt.fields,
			}
			got := s.Get()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricStorerImpl_Set(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]models.Metric
		args   map[string]models.Metric
		want   map[string]models.Metric
	}{
		{
			name:   "Set on empty store",
			fields: map[string]models.Metric{},
			args: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
			},
			want: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
			},
		},
		{
			name: "Set on existing store",
			fields: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
			},
			args: map[string]models.Metric{
				"metric2": {ID: "metric2", MType: "counter", Delta: new(int64)},
			},
			want: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: "gauge", Value: new(float64)},
				"metric2": {ID: "metric2", MType: "counter", Delta: new(int64)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricStorerImpl{
				metrics: tt.fields,
			}
			s.Set(tt.args)
			if got := s.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("New returns valid MetricStorerImpl instance", func(t *testing.T) {
		got := New()
		if got == nil {
			t.Errorf("New() = nil, want instance of MetricStorer")
		}
		if _, ok := got.(*MetricStorerImpl); !ok {
			t.Errorf("New() = %T, want *MetricStorerImpl", got)
		}
	})
}
