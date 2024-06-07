package main

import (
	"net/http"

	"github.com/VOTONO/go-metrics/internal/handlers"
	"github.com/VOTONO/go-metrics/internal/storage"
)

func main() {
	memStorage := storage.New(nil)

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateHandler(memStorage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
