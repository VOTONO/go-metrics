package workers

import (
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/models"
)

func BenchmarkRead(b *testing.B) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	readResultChannel := make(chan []models.Metric, 1)

	readWorker := NewReadWorker(
		sugaredLogger,
		readResultChannel,
		10,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		readWorker.readRuntimeMetrics()
	}
}
