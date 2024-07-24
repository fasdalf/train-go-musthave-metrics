package main

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/jsonofflinestorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	c := config.GetConfig()
	memStorage := metricstorage.NewMemStorageWithSave()
	fileStorage := jsonofflinestorage.NewJSONFileStorage(memStorage, c.StorageFileName, c.StorageFileRestore, c.StorageFileStoreInterval)
	if err := fileStorage.Restore(); err != nil {
		panic(err)
	}

	engine := server.NewHTTPEngine(memStorage, fileStorage)
	srv := &http.Server{
		Addr:    c.Addr,
		Handler: engine,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		slog.Info("receive interrupt signal")
		if err := srv.Close(); err != nil {
			slog.Error("Server Close:", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			slog.Info("Server closed under request")
		} else {
			panic(err)
		}
	}

	fileStorage.Save()
}
