package main

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/config"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/handlers"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"time"
)

func main() {
	collectInterval := time.Duration(config.GetConfig().PollInterval) * time.Second
	sendInterval := time.Duration(config.GetConfig().ReportInterval) * time.Second
	address := config.GetConfig().Addr
	memStorage := metricstorage.NewMemStorage()

	collectTimeout := time.Duration(0)
	sendTimeout := time.Duration(0)
	for {
		if collectTimeout <= time.Duration(0) {
			collectTimeout = collectInterval
			handlers.CollectMetrics(memStorage)
		}
		if sendTimeout <= time.Duration(0) {
			sendTimeout = sendInterval
			handlers.SendMetrics(memStorage, address)
		}

		sleepTime := min(collectTimeout, sendTimeout)
		collectTimeout -= sleepTime
		sendTimeout -= sleepTime
		fmt.Println("sleeping for", sleepTime)
		time.Sleep(sleepTime)
	}
}
