package main

import (
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server/router"
	"github.com/VOTONO/go-metrics/internal/server/storage"
)

func main() {
	config := getConfig()
	stor := storage.New(nil)
	router := router.Router(stor)

	err := http.ListenAndServe(config.address, router)
	if err != nil {
		panic(err)
	}
}
