package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"net/http"
	"time"
)

func BatchUpdateHandler(s repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metrics []models.Metric
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 1000*time.Second)
		defer cancel()

		er := s.StoreSlice(ctx, metrics)
		if er != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
}
