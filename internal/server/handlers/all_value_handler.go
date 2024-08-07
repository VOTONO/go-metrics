package handlers

import (
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"go.uber.org/zap"
	"net/http"
)

func AllValueHandler(s repo.MetricStorer, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metrics, err := s.All()

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		htmlContent, err := helpers.MetricsToHTML(metrics, logger)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, htmlContent)
	}
}
