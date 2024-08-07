package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VOTONO/go-metrics/internal/agent/monitor"
	"github.com/VOTONO/go-metrics/internal/agent/network"
	"github.com/VOTONO/go-metrics/internal/agent/storage"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()
	config := getConfig()

	readTicker := time.NewTicker(time.Duration(config.pollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.reportInterval) * time.Second)

	net := network.New(&http.Client{}, config.address, sugar)
	stor := storage.New()

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
