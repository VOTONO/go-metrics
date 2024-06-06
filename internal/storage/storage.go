package storage

import (
	"fmt"
)

type MetricStorage interface {
	Replace(name string, value interface{}) error
	Increment(name string, value interface{}) error
}

type MemStorage struct {
	metrics map[string]interface{}
}

func New() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func (m *MemStorage) Replace(name string, value interface{}) error {
	switch v := value.(type) {
	case int64:
		m.metrics[name] = v
		fmt.Println("Replaced metric", name, "with value", v)
	case float64:
		m.metrics[name] = v
		fmt.Println("Replaced metric", name, "with value", v)
	default:
		return fmt.Errorf("unsupported type %T for metric value", value)
	}
	return nil
}

func (m *MemStorage) Increment(name string, value interface{}) error {
	switch v := value.(type) {
	case int64:
		if val, ok := m.metrics[name].(int64); ok {
			m.metrics[name] = val + v
			fmt.Println("Incremented metric", name, "by", v)
		} else {
			return fmt.Errorf("metric %s value is not int64", name)
		}
	case float64:
		if val, ok := m.metrics[name].(float64); ok {
			m.metrics[name] = val + v
			fmt.Println("Incremented metric", name, "by", v)
		} else {
			return fmt.Errorf("metric %s value is not float64", name)
		}
	default:
		return fmt.Errorf("unsupported type %T for metric value", value)
	}
	return nil
}
