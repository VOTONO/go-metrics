package repo

import (
	"context"
	"go.uber.org/zap"
	"time"
)

func StartWriting(ctx context.Context, storer MetricStorer, logger *zap.SugaredLogger, storeInterval int, filePath string) {
	if storeInterval <= 0 || filePath == "" {
		return
	}

	storeTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)

	go func() {
		defer storeTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				Write(filePath, storer.All(), logger)
				return
			case <-storeTicker.C:
				Write(filePath, storer.All(), logger)
			}
		}
	}()
}
