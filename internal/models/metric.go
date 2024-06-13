package models

import (
	"fmt"
	"strconv"
)

type Metric struct {
	Name  string
	Type  string
	Value string
}

// Constructor for Metric with validation
func New(name, metricType, value string) (Metric, error) {
	if name == "" {
		return Metric{}, fmt.Errorf("bad name")
	}

	if metricType != "gauge" && metricType != "counter" {
		return Metric{}, fmt.Errorf("bad type")
	}

	if metricType == "gauge" {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return Metric{}, fmt.Errorf("bad value")
		}
	}

	if metricType == "counter" {
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return Metric{}, fmt.Errorf("bad value")
		}
	}

	return Metric{Name: name, Type: metricType, Value: value}, nil
}

// Add another Metric's value to this one. If they are of the same type and compatible, returns new metric.
func (m Metric) Add(other Metric) (Metric, error) {
	if m.Type != other.Type {
		return Metric{}, fmt.Errorf("type mismatch: cannot add metric of type %s to type %s", other.Type, m.Type)
	}

	newMetric := Metric{Name: m.Name, Type: m.Type}

	switch m.Type {
	case "gauge":
		v, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid value for gauge: %s", m.Value)
		}
		ov, err := strconv.ParseInt(other.Value, 10, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid value for gauge: %s", other.Value)
		}
		newMetric.Value = strconv.FormatInt(v+ov, 10)
	case "counter":
		v, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid value for counter: %s", m.Value)
		}
		ov, err := strconv.ParseFloat(other.Value, 64)
		if err != nil {
			return Metric{}, fmt.Errorf("invalid value for counter: %s", other.Value)
		}
		sum := v + ov
		newMetric.Value = strconv.FormatFloat(sum, 'G', -1, 64)
	default:
		return Metric{}, fmt.Errorf("unsupported type %s for metric value", m.Type)
	}

	return newMetric, nil
}
