package handlers

import (
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"net/http"
)

func AllValueHandler(s repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metrics, err := s.All()

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, "<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>")
		for key, metric := range metrics {
			value, err := helpers.ExtractValue(metric)

			if err != nil {
				http.Error(res, "Invalid metric value", http.StatusInternalServerError)
			}

			fmt.Fprintf(res, "<tr><td>%s</td><td>%v</td></tr>", key, value)
		}
		fmt.Fprintln(res, "</table></body></html>")
	}
}
