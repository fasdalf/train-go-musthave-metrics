package main

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/config"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/handlers"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"log/slog"
	"time"
)

func main() {
	collectInterval := time.Duration(config.GetConfig().PollInterval) * time.Second
	sendInterval := time.Duration(config.GetConfig().ReportInterval) * time.Second
	address := config.GetConfig().Addr
	memStorage := metricstorage.NewMemStorage()
	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})

	time.AfterFunc(100*time.Millisecond, func() {
		go sendMetricsRoutine(memStorage, address, sendInterval, retryer)
	})
	for {
		handlers.CollectMetrics(memStorage)
		slog.Info(`collector sleeping`, `delay`, collectInterval)
		time.Sleep(collectInterval)
	}
}

func sendMetricsRoutine(storage handlers.Storage, address string, sendInterval time.Duration, retryer handlers.Retryer) {
	for {
		handlers.SendMetrics(storage, address, retryer)

		slog.Info(`sender sleeping`, `delay`, sendInterval)
		time.Sleep(sendInterval)
	}
}
