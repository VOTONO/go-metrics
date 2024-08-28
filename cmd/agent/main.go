package main

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/helpers"
	"github.com/VOTONO/go-metrics/internal/agent/workers"
	"github.com/VOTONO/go-metrics/internal/models"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	config := getConfig()

	sugaredLogger.Infow("starting agent",
		"address", config.address,
		"pollInterval", config.pollInterval,
		"reportInterval", config.reportInterval,
		"secretKey", config.secretKey,
	)

	stopChannel := helpers.CreateSystemStopChannel()
	readResultChannel := make(chan []models.Metric, 1)

	readWorker := workers.NewReadWorker(
		sugaredLogger,
		readResultChannel,
		config.pollInterval,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	sendWorker := workers.NewSendWorker(
		client,
		sugaredLogger,
		config.reportInterval,
		readResultChannel,
		config.rateLimit,
		config.address,
		config.secretKey,
	)

	go func() {
		readWorker.Start()
	}()

	go func() {
		sendWorker.Start()
	}()

	<-stopChannel
	logger.Info("stopping agent")
	readWorker.Stop()
	sendWorker.Stop()
}
