package workers

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/models"
)

func ExampleReadWorker_Start() {
	// Init logger for worker.
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	// Init read result channel.
	readResultChannel := make(chan []models.Metric, 1)

	// Catch syscall to stop worker.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)
	defer stop()

	// Init worker.
	worker := NewReadWorker(
		sugaredLogger,
		readResultChannel, // Metrics will be sent in this channel.
		2,                 // Interval in seconds of metrics reading.
	)

	// Call Start() function in goroutine.
	go func() {
		worker.Start()
	}()

	select {
	// Metrics will be sent in this channel periodically.
	case readResult := <-readResultChannel:
		fmt.Println(readResult)
	case <-ctx.Done():
		// Stop worker gracefully.
		worker.Stop()
	}

}
