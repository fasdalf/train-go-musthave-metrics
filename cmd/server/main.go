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
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
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

	var db handlers.Pingable

	switch true {
	case c.StorageDBDSN != "":
		slog.Info("initializing database connection", "DATABASE_DSN", c.StorageDBDSN)
		var err error
		pgx, err := sql.Open("pgx", c.StorageDBDSN)
		if err != nil {
			slog.Error("can not connect to DB", "error", err)

			panic(err)
		}

		defer pgx.Close()

		dbStorage, err := metricstorage.NewDBStorage(pgx, context.Background())
		if err != nil {
			slog.Error("can not init DB", "error", err)
			panic(err)
		}
		metricStorage = metricstorage.NewSavableModelStorage(dbStorage)
		db = pgx
	case c.StorageFileName != "":
		slog.Info("initializing in-mem and in-file storage")
		dirtyStorage := metricstorage.NewDirtyStorage(metricstorage.NewMemStorage())
		modelStorage := metricstorage.NewSavableModelStorage(dirtyStorage)

		slog.Info("initializing file storage")
		fileSaver := jsonofflinestorage.NewJSONFileStorage(modelStorage, c.StorageFileName, c.StorageFileRestore, c.StorageFileStoreInterval, dirtyStorage.SavedChan, dirtyStorage.Clear)
		if err := fileSaver.Start(ctx, wg); err != nil {
			slog.Error("can not init file storage", "error", err)
			panic(err)
		}
		metricStorage = modelStorage
	default:
		slog.Info("initializing in-mem only storage")
		metricStorage = metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())
	}

	slog.Debug("initializing http router")
	engine := server.NewRoutingEngine(metricStorage, db, retryer, c.HashKey)
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
