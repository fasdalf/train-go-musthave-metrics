package main

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/config"
)

func main() {
	memStorage := metricstorage.NewMemStorage()
	httpServer := server.NewHTTPEngine(memStorage)
	err := httpServer.Run(config.GetConfig().Addr)
	if err != nil {
		panic(err)
	}
}
