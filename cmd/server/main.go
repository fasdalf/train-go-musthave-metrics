package main

import (
	"context"
	"database/sql"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/jsonofflinestorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := config.GetConfig()
	slog.Info("initializing server")

	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	metricStorage := (*metricstorage.SavableModelStorage)(nil)

	db := (*sql.DB)(nil)
	if c.StorageDBDSN != "" {
		slog.Info("initializing database connection", "DATABASE_DSN", c.StorageDBDSN)
		err := error(nil)
		db, err = sql.Open("pgx", c.StorageDBDSN)
		if err != nil {
			slog.Error("can not connect to DB", "error", err)

			panic(err)
		}

		defer db.Close()

		dbStorage, err := metricstorage.NewDBStorage(db, context.Background())
		if err != nil {
			slog.Error("can not init DB", "error", err)
			panic(err)
		}
		metricStorage = metricstorage.NewSavableModelStorage(dbStorage)
	}

	if metricStorage == nil {
		slog.Info("initializing in-mem storage")
		metricStorage = metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())
	}

	savedChan := (jsonofflinestorage.SavedChan)(nil)
	if c.StorageFileName != "" && db == nil {
		slog.Info("initializing file storage")
		fileStorageService := jsonofflinestorage.NewJSONFileStorage(metricStorage, c.StorageFileName, c.StorageFileRestore, c.StorageFileStoreInterval)
		if err := fileStorageService.Restore(); err != nil {
			slog.Error("can not read file storage", "error", err)
			panic(err)
		}
		defer fileStorageService.Save()

		savedChan = make(jsonofflinestorage.SavedChan)
		go func() {
			err := fileStorageService.SaverRoutine(savedChan)
			if err != nil {
				slog.Error("async saver failed", "error", err)
			}
		}()
	}

	slog.Debug("initializing http router")
	engine := server.NewRoutingEngine(metricStorage, savedChan, db, retryer)
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
