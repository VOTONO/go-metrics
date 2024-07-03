package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/storage"
)

func UpdateHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metric
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		valid := validateMetric(metric)
		if !valid {
			http.Error(res, "invalid metric", http.StatusBadRequest)
			return
		}

		stored, err := s.Store(metric)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		out, err := json.Marshal(stored)
		if err != nil {
			log.Fatal(err)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}

func ValueHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metric
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		value, found := s.Get(metric.ID)

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}

		out, err := json.Marshal(value)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}

func AllValueHandler(memStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metrics := memStorage.All()

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, "<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>")
		for key, value := range metrics {
			fmt.Fprintf(res, "<tr><td>%s</td><td>%v</td></tr>", key, value)
		}
		fmt.Fprintln(res, "</table></body></html>")
	}
}

func validateMetric(m models.Metric) bool {
	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return false
		}
		return true
	case "counter":
		if m.Delta == nil {
			return false
		}
		return true
	default:
		return false
	}
}
