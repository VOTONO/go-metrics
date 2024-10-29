package handlers

import (
	"database/sql"
	"net/http"
)

func Ping(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, _ *http.Request) {
		if db == nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		err := db.Ping()

		if err != nil {
			http.Error(res, "No connection to db", http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
