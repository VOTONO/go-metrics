package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VOTONO/go-metrics/internal/agent/logic"
	"github.com/VOTONO/go-metrics/internal/agent/network"
	"github.com/VOTONO/go-metrics/internal/agent/repo"
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	metricReader := logic.NewMetricReaderImpl()
	metricSender := network.New(client, config.address, sugar)
	metricStorage := repo.New()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			readTicker.Stop()
			sendTicker.Stop()
			return
		case <-readTicker.C:
			metrics := metricReader.Read()
			metricStorage.Set(metrics)
		case <-sendTicker.C:
			metrics := metricStorage.Get()
			metricSender.Send(metrics)
		}
	}
}
