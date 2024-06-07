package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type Metric struct {
	Name  string
	Type  string
	Value interface{}
}

var metrics = map[string]Metric{
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

func main() {
	updateTicker := time.NewTicker(2 * time.Second)
	sendTicker := time.NewTicker(10 * time.Second)

	client := &http.Client{}

	go func() {
		for range updateTicker.C {
			updateMetricsWith(runtime.MemStats{})
		}
	}()

	go func() {
		for range sendTicker.C {
			sendMetrics(client)
		}
	}()

	select {}
}

func sendMetrics(client *http.Client) {
	for _, metric := range metrics {

		req, err := requestFor(metric)
		if err != nil {
			fmt.Printf("Error creating HTTP request for metric %s: %v\n", metric, err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending HTTP request for metric %s: %v\n", metric, err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for metric %s failed with status code %d\n", metric, resp.StatusCode)
		}

		resp.Body.Close()
	}
}

func requestFor(metric Metric) (*http.Request, error) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metric.Name, metric.Type, metric.Value)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	return req, nil
}

func updateMetricsWith(memStats runtime.MemStats) {
	if pollCountMetric, found := metrics["PollCount"]; found {
		pollCountMetric.Value = pollCountMetric.Value.(int) + 1
		metrics["PollCount"] = pollCountMetric
	}
	random := rand.Intn(100)
	metrics["RandomValue"] = Metric{Name: "RandomValue", Type: "gauge", Value: random}

	runtime.ReadMemStats(&memStats)

	val := reflect.ValueOf(memStats)
	typ := reflect.TypeOf(memStats)

	for i := 0; i < val.NumField(); i++ {
		name := typ.Field(i).Name
		if _, found := metrics[name]; found {
			value := val.Field(i).Interface()
			metrics[name] = Metric{Name: name, Type: "gauge", Value: value}
		}
	}
}
