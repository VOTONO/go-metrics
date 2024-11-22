package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/helpers"
	"github.com/VOTONO/go-metrics/internal/agent/workers"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	config := getConfig()

	sugaredLogger.Infow(
		"Ldflags",
		"Build version", buildVersion,
		"Build date", buildDate,
		"Build commit", buildCommit,
	)
	sugaredLogger.Infow("starting agent",
		"Address", config.Address,
		"PollInterval", config.PollInterval,
		"ReportInterval", config.ReportInterval,
		"SecretKey", config.SecretKey,
		"PublicKeyPath", config.PublicKeyPath,
	)

	stopChannel := helpers.CreateSystemStopChannel()

	readWorker := workers.NewReadWorker(
		sugaredLogger,
		config.PollInterval,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// If public key provided, add TLS to client
	if config.PublicKeyPath != "" {
		serverCert, err := os.ReadFile(config.PublicKeyPath)
		if err != nil {
			sugaredLogger.Fatalf("failed to read server certificate: %v", err)
		}

		// Create a certificate pool and add the server's certificate
		certPool := x509.NewCertPool()
		certAdded := certPool.AppendCertsFromPEM(serverCert)
		if !certAdded {
			sugaredLogger.Fatal("failed to append server certificate to cert pool")
		}

		// Configure the TLS settings for the client
		tlsConfig := &tls.Config{
			RootCAs: certPool,
		}

		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	sendWorker := workers.NewSendWorker(
		client,
		sugaredLogger,
		config.ReportInterval,
		readWorker.ResultChannel,
		config.RateLimit,
		config.Address,
		config.SecretKey,
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
