package main

import (
	"context"
	"database/sql"
	"errors"
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
	"sync"
	"time"
)

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	c := config.GetConfig()
	slog.Info("initializing server")

	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	var metricStorage *metricstorage.SavableModelStorage

	var db *sql.DB = nil
	var savedChan jsonofflinestorage.SavedChan = nil

	switch true {
	case c.StorageDBDSN != "":
		slog.Info("initializing database connection", "DATABASE_DSN", c.StorageDBDSN)
		var err error
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
	case c.StorageFileName != "":
		slog.Info("initializing in-mem and in-file storage")
		metricStorage = metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())

		slog.Info("initializing file storage")
		fileStorageService := jsonofflinestorage.NewJSONFileStorage(metricStorage, c.StorageFileName, c.StorageFileRestore, c.StorageFileStoreInterval)
		if err := fileStorageService.Restore(); err != nil {
			slog.Error("can not read file storage", "error", err)
			panic(err)
		}

		if c.StorageFileStoreInterval > 0 {
			savedChan = make(jsonofflinestorage.SavedChan)
		}

		wg.Add(1)
		go func() {
			err := fileStorageService.SaveMetrics(ctx, savedChan, wg)
			if err != nil {
				slog.Error("async saver failed", "error", err)
			}
		}()
	default:
		slog.Info("initializing in-mem only storage")
		metricStorage = metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())
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
		ctxCancel()
		if err := srv.Close(); err != nil {
			slog.Error("Server close error:", "error", err)
		}
	}()

	slog.Info("starting http server", "address", c.Addr)
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("Server closed by interrupt signal")
			slog.Info("wait for bg processes")
			wg.Wait()
		} else {
			slog.Error("server not started or stopped with error", "error", err)
			panic(err)
		}
	}
}
