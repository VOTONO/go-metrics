package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"net/http"
)

func ValueHandlerJSON(s repo.MetricStorer) http.HandlerFunc {
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
