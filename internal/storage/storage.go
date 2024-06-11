package storage

import (
	"fmt"
)

type MetricStorage interface {
	Replace(name string, value interface{}) error
	Increment(name string, value interface{}) error
	Get(name string) interface{}
	GetAll() map[string]interface{}
}

type MemStorage struct {
	metrics map[string]interface{}
}

func New(storage map[string]interface{}) *MemStorage {
	if storage == nil {
		storage = make(map[string]interface{})
	}
	return &MemStorage{
		metrics: storage,
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
		if val, ok := m.metrics[name]; ok {
			if existingVal, ok := val.(int64); ok {
				m.metrics[name] = existingVal + v
				fmt.Println("Incremented metric", name, "by", v)
			} else {
				return fmt.Errorf("metric %s value is not int64", name)
			}
		} else {
			m.metrics[name] = v
			fmt.Println("Added new int64 metric", name, "with value", v)
		}
	case float64:
		if val, ok := m.metrics[name]; ok {
			if existingVal, ok := val.(float64); ok {
				m.metrics[name] = existingVal + v
				fmt.Println("Incremented metric", name, "by", v)
			} else {
				return fmt.Errorf("metric %s value is not float64", name)
			}
		} else {
			m.metrics[name] = v
			fmt.Println("Added new float64 metric", name, "with value", v)
		}
	default:
		return fmt.Errorf("unsupported type %T for metric value", value)
	}
	return nil
}

func (m *MemStorage) Get(name string) interface{} {
	return m.metrics[name]
}

func (m *MemStorage) GetAll() map[string]interface{} {
	return m.metrics
}
