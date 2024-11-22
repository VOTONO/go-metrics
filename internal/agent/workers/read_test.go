package workers

import (
	"log"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkRead(b *testing.B) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	readWorker := NewReadWorker(
		sugaredLogger,
		10,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		readWorker.readRuntimeMetrics()
	}
}
