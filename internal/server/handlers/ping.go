package handlers

import (
	"database/sql"
	"net/http"
)

func Ping(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := db.Ping()

		if err != nil {
			http.Error(res, "No connection to db", http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
