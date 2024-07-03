package monitor

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/VOTONO/go-metrics/internal/models"
)

var count = int64(0)

var metrics = map[string]models.Metric{
	"PollCount":     {},
	"RandomValue":   {},
	"Alloc":         {},
	"BuckHashSys":   {},
	"Frees":         {},
	"GCCPUFraction": {},
	"GCSys":         {},
	"HeapAlloc":     {},
	"HeapIdle":      {},
	"HeapInuse":     {},
	"HeapObjects":   {},
	"HeapReleased":  {},
	"HeapSys":       {},
	"LastGC":        {},
	"Lookups":       {},
	"MCacheInuse":   {},
	"MCacheSys":     {},
	"MSpanInuse":    {},
	"MSpanSys":      {},
	"Mallocs":       {},
	"NextGC":        {},
	"NumForcedGC":   {},
	"NumGC":         {},
	"OtherSys":      {},
	"PauseTotalNs":  {},
	"StackInuse":    {},
	"StackSys":      {},
	"Sys":           {},
	"TotalAlloc":    {},
}

var stats runtime.MemStats

func Read() map[string]models.Metric {
	count++
	metrics["PollCount"] = models.Metric{
		ID:    "PollCount",
		MType: "counter",
		Delta: &count,
	}
	random := rand.Float64()
	metrics["RandomValue"] = models.Metric{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &random,
	}

	runtime.ReadMemStats(&stats)

	val := reflect.ValueOf(stats)
	typ := reflect.TypeOf(stats)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name := typ.Field(i).Name

		if _, found := metrics[name]; found {
			var value float64

			switch field.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value = float64(field.Uint())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value = float64(field.Int())
			case reflect.Float32, reflect.Float64:
				value = field.Float()
			default:
				continue
			}

			metrics[name] = models.Metric{
				ID:    name,
				MType: "gauge",
				Value: &value,
			}
		}
	}

	return metrics
}
