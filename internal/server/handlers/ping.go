package handlers

import (
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"net/http"
)

func Ping(s repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := s.Ping()

		if err != nil {
			http.Error(res, "No connection to db", http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
