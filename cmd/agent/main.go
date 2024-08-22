package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/logic"
	"github.com/VOTONO/go-metrics/internal/agent/network"
	"github.com/VOTONO/go-metrics/internal/agent/repo"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := *logger.Sugar()
	config := getConfig()

	sugaredLogger.Infow("starting agent",
		"address", config.address,
		"pollInterval", config.pollInterval,
		"reportInterval", config.reportInterval,
		"secretKey", config.secretKey,
	)

	readTicker := time.NewTicker(time.Duration(config.pollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.reportInterval) * time.Second)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	metricReader := logic.NewMetricReaderImpl()
	metricSender := network.New(client, config.address, &sugaredLogger)
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
			err := metricSender.Send(metrics, config.secretKey)
			if err != nil {
				sugaredLogger.Errorw("send metrics failed", "error", err.Error())
			}
		}
	}
}
