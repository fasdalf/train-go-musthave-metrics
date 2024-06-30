package main

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/client/handlers"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"time"
)

func main() {
	const collectInterval = 2 * time.Second
	const sendInterval = 10 * time.Second
	memStorage := metricstorage.NewMemStorage()

	collectTimeout := time.Duration(0)
	sendTimeout := time.Duration(0)
	for {
		// Do i need https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i to stop eternal loop?
		// No, i don't
		if collectTimeout <= time.Duration(0) {
			collectTimeout = collectInterval
			handlers.CollectMetrics(memStorage)
		}
		if sendTimeout <= time.Duration(0) {
			sendTimeout = sendInterval
			handlers.SendMetrics(memStorage)
		}

		sleepTime := min(collectTimeout, sendTimeout)
		collectTimeout -= sleepTime
		sendTimeout -= sleepTime
		fmt.Println("sleeping for", sleepTime)
		time.Sleep(sleepTime)
	}
}
