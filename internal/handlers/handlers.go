package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/VOTONO/go-metrics/internal/storage"
)

func UpdateHandler(memStorage storage.MetricStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		}

		path := strings.TrimLeft(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) != 4 {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metricType := pathParts[1]
		metricName := pathParts[2]
		metricValue := pathParts[3]

		switch metricType {
		case "gauge":
			floatValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(res, "Invalid value", http.StatusBadRequest)
				return
			}

			if err := memStorage.Replace(metricName, floatValue); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		case "counter":
			intValue, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(res, "Invalid value", http.StatusBadRequest)
				return
			}

			if err := memStorage.Increment(metricName, intValue); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(res, "Invalid metric type", http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}

func ValueHandler(memStorage storage.MetricStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		}

		path := strings.TrimLeft(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) != 3 {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metricName := pathParts[2]

		value := memStorage.Get(metricName)

		if value == nil {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(fmt.Sprintf("%v", value)))
	}
}

func AllValueHandler(memStorage storage.MetricStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		}

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metrics := memStorage.GetAll()

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, "<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>")
		for key, value := range metrics {
			fmt.Fprintf(res, "<tr><td>%s</td><td>%v</td></tr>", key, value)
		}
		fmt.Fprintln(res, "</table></body></html>")
	}
}
