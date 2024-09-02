package helpers

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateSystemStopChannel() chan os.Signal {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	return stopChannel
}

func CreateTicker(seconds int) *time.Ticker {
	duration := time.Duration(seconds) * time.Second
	return time.NewTicker(duration)
}
