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
	"sync"
	"time"
)

func main() {
	cfg := config.GetConfig()
	collectInterval := time.Duration(cfg.PollInterval) * time.Second
	sendInterval := time.Duration(cfg.ReportInterval) * time.Second
	address := cfg.Addr
	memStorage := metricstorage.NewMemStorageMuted()
	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go handlers.SendMetricsLoop(ctx, wg, memStorage, address, sendInterval, retryer, cfg.HashKey, cfg.RateLimit)
	go handlers.Collect(handlers.CollectMetrics, ctx, wg, memStorage, collectInterval)
	go handlers.Collect(handlers.CollectGopsutilMetrics, ctx, wg, memStorage, collectInterval)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("interrupt signal received")
	signal.Stop(quit)
	cancel()
	slog.Info("attempting graceful shutdown")
	wg.Wait()
}
