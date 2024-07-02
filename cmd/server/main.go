package main

import (
	"log"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server/router"
	"github.com/VOTONO/go-metrics/internal/server/storage"

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
	stor := storage.New(nil)
	router := router.Router(stor, sugar)

	sugar.Infow(
		"Starting server",
		"address", config.address,
	)
	er := http.ListenAndServe(config.address, router)
	if er != nil {
		sugar.Errorw(
			"Fail start server",
			"error", er,
		)
	}
}
