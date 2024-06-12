package main

import (
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server"
	"github.com/VOTONO/go-metrics/internal/server/storage"
)

func main() {
	config := getConfig()

	stor := storage.New(nil)

	httpServer := server.New(stor)

	err := http.ListenAndServe(config.address, httpServer)
	if err != nil {
		panic(err)
	}
}
