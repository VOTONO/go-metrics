package main

import (
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server"
	"github.com/VOTONO/go-metrics/internal/storage"
)

func main() {
	memStorage := storage.New(nil)
	httpServer := server.New(memStorage)

	err := http.ListenAndServe(":8080", httpServer)
	if err != nil {
		panic(err)
	}
}
