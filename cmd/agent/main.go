package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VOTONO/go-metrics/internal/agent/monitor"
	"github.com/VOTONO/go-metrics/internal/agent/network"
	"github.com/VOTONO/go-metrics/internal/agent/storage"
	"github.com/VOTONO/go-metrics/internal/models"
)

func main() {
	config := getConfig()

	readTicker := time.NewTicker(time.Duration(config.pollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.reportInterval) * time.Second)

	net := network.New(&http.Client{}, config.address)
	stor := storage.New(map[string]models.Metric{})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			readTicker.Stop()
			sendTicker.Stop()
			return
		case <-readTicker.C:
			metrics := monitor.Read()
			stor.Set(metrics)
		case <-sendTicker.C:
			metrics := stor.Get()
			net.Send(metrics)
		}
	}
}
