package monitor

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/VOTONO/go-metrics/internal/models"
)

var count = 0

var metrics = map[string]models.Metric{
	"PollCount":     {Name: "PollCount", Type: "counter", Value: 0},
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
	metrics["PollCount"] = models.Metric{Name: "RandomValue", Type: "gauge", Value: count + 1}
	metrics["RandomValue"] = models.Metric{Name: "RandomValue", Type: "gauge", Value: rand.Intn(100)}

	runtime.ReadMemStats(&stats)

	val := reflect.ValueOf(stats)
	typ := reflect.TypeOf(stats)

	for i := 0; i < val.NumField(); i++ {
		name := typ.Field(i).Name
		if _, found := metrics[name]; found {
			value := val.Field(i).Interface()
			metrics[name] = models.Metric{Name: name, Type: "gauge", Value: value}
		}
	}

	return metrics
}
