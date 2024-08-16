package main

import (
	"context"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/config"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/handlers"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func main() {
	collectInterval := time.Duration(config.GetConfig().PollInterval) * time.Second
	sendInterval := time.Duration(config.GetConfig().ReportInterval) * time.Second
	address := config.GetConfig().Addr
	memStorage := metricstorage.NewMemStorageMuted()
	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())

	go sendMetrics(ctx, memStorage, address, sendInterval, retryer)
	go collectMetrics(ctx, memStorage, collectInterval)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("interrupt signal received")
	signal.Stop(quit)
	cancel()
	slog.Info("attempting graceful shutdown")
	time.Sleep(time.Second)
}

func collectMetrics(ctx context.Context, storage handlers.Storage, collectInterval time.Duration) {
main:
	for {
		select {
		case <-ctx.Done():
			break main
		default:
		}
		handlers.CollectMetrics(storage)
		slog.Info(`collector sleeping`, `delay`, collectInterval)
		time.Sleep(collectInterval)
	}
}

func sendMetrics(
	ctx context.Context,
	storage handlers.Storage,
	address string,
	sendInterval time.Duration,
	retryer handlers.Retryer,
	key string,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		slog.Info(`sender sleeping`, `delay`, sendInterval)
		time.Sleep(sendInterval)
		handlers.SendMetrics(ctx, storage, address, retryer, key)
	}
}
