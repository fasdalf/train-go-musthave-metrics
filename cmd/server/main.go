package main

import (
	"database/sql"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/jsonofflinestorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	c := config.GetConfig()
	slog.Info("initializing mem storage")
	memStorage := metricstorage.NewMemStorageWithSave()

	fileStorage := (handlers.FileStorage)(nil)
	if c.StorageFileName != "" {
		slog.Info("initializing file storage")
		fileStorageService := jsonofflinestorage.NewJSONFileStorage(memStorage, c.StorageFileName, c.StorageFileRestore, c.StorageFileStoreInterval)
		if err := fileStorageService.Restore(); err != nil {
			slog.Error("can not read file storage", "error", err)
			panic(err)
		}
		defer fileStorageService.Save()
		fileStorage = fileStorageService
	}

	db := (*sql.DB)(nil)
	if c.StorageDBDSN != "" {
		err := error(nil)
		db, err = sql.Open("pgx", c.StorageDBDSN)
		if err != nil {
			panic(err)
		}

		defer db.Close()
	}

	slog.Debug("initializing http router")
	engine := server.NewRoutingEngine(memStorage, fileStorage, db)
	srv := &http.Server{
		Addr:    c.Addr,
		Handler: engine,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		slog.Info("interrupt signal received")
		signal.Stop(quit)
		if err := srv.Close(); err != nil {
			slog.Error("Server close error:", "error", err)
		}
	}()

	slog.Info("starting http server", "address", c.Addr)
	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			slog.Info("Server closed by interrupt signal")
		} else {
			slog.Error("server not started or stopped with error", "error", err)
			panic(err)
		}
	}
}
