package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server"
	"github.com/VOTONO/go-metrics/internal/storage"
)

func main() {
	defaultAddress := &NetAddress{"localhost", 8080}
	addr := defaultAddress

	flag.Var(addr, "a", "Net address host:port")

	flag.Parse()

	fmt.Println("Address:", addr.String())

	memStorage := storage.New(nil)

	httpServer := server.New(memStorage)

	err := http.ListenAndServe(addr.String(), httpServer)
	if err != nil {
		panic(err)
	}
}
