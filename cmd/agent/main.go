package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/helpers"
	"github.com/VOTONO/go-metrics/internal/agent/workers"
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

	readWorker := workers.NewReadWorker(
		sugaredLogger,
		config.pollInterval,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	sendWorker := workers.NewSendWorker(
		client,
		sugaredLogger,
		config.reportInterval,
		readWorker.ResultChannel,
		config.rateLimit,
		config.address,
		config.secretKey,
	)

	go func() {
		err := http.ListenAndServe(":9191", nil)
		if err != nil {
			sugaredLogger.Errorw("Fail start agent", "error", err)
		}
	}()

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
