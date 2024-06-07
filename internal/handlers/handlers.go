package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/VOTONO/go-metrics/internal/storage"
)

func UpdateHandler(memStorage storage.MetricStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimLeft(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) != 4 {
			http.Error(res, "Invalid request format", http.StatusNotFound)
			return
		}

		metricType := pathParts[1]
		metricName := pathParts[2]
		metricValue := pathParts[3]

		switch metricType {
		case "gauge":
			floatValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(res, "Invalid value for gauge, must be float64", http.StatusBadRequest)
				return
			}

			if err := memStorage.Replace(metricName, floatValue); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		case "counter":
			intValue, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(res, "Invalid value for counter, must be int64", http.StatusBadRequest)
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
	}
}
