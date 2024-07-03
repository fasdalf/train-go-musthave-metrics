package main

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
)

func main() {
	memStorage := metricstorage.NewMemStorage()
	httpServer := server.NewHTTPServer(memStorage)
	err := httpServer.Run(`:8080`)
	if err != nil {
		panic(err)
	}
}
