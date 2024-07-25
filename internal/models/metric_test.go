package models

import (
	"reflect"
	"testing"
)

func TestNewMetric(t *testing.T) {
	float64Ptr := func(v float64) *float64 { return &v }
	int64Ptr := func(v int64) *int64 { return &v }

	type args struct {
		id         string
		metricType string
		value      string
	}
	tests := []struct {
		name    string
		args    args
		want    Metric
		wantErr bool
	}{
		{
			name: "Valid gauge metric",
			args: args{
				id:         "testGauge",
				metricType: "gauge",
				value:      "123.45",
			},
			want: Metric{
				ID:    "testGauge",
				MType: "gauge",
				Value: float64Ptr(123.45),
			},
			wantErr: false,
		},
		{
			name: "Valid counter metric",
			args: args{
				id:         "testCounter",
				metricType: "counter",
				value:      "123",
			},
			want: Metric{
				ID:    "testCounter",
				MType: "counter",
				Delta: int64Ptr(123),
			},
			wantErr: false,
		},
		{
			name: "Invalid gauge value",
			args: args{
				id:         "invalidGauge",
				metricType: "gauge",
				value:      "abc",
			},
			want:    Metric{},
			wantErr: true,
		},
		{
			name: "Invalid counter value",
			args: args{
				id:         "invalidCounter",
				metricType: "counter",
				value:      "abc",
			},
			want:    Metric{},
			wantErr: true,
		},
		{
			name: "Invalid metric type",
			args: args{
				id:         "invalidType",
				metricType: "invalid",
				value:      "123",
			},
			want:    Metric{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.id, tt.args.metricType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}
