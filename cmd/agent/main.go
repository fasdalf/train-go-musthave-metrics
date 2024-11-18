package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/config"
	"github.com/fasdalf/train-go-musthave-metrics/internal/agent/handlers"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/printbuild"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	const pprofHTTPAddr = ":8092"

	bd := &printbuild.Data{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}
	bd.Print()

	cfg := config.GetConfig()
	collectInterval := time.Duration(cfg.PollInterval) * time.Second
	sendInterval := time.Duration(cfg.ReportInterval) * time.Second
	address := cfg.Addr
	memStorage := metricstorage.NewMemStorageMuted()
	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go handlers.SendMetricsLoop(ctx, wg, memStorage, address, sendInterval, retryer, handlers.NewNetHTTPPoster(), cfg.HashKey, cfg.RSAKey, cfg.RateLimit)
	go handlers.Collect(handlers.CollectMetrics, ctx, wg, memStorage, collectInterval)
	go handlers.Collect(handlers.CollectGopsutilMetrics, ctx, wg, memStorage, collectInterval)
	// for "net/http/pprof"
	go http.ListenAndServe(pprofHTTPAddr, nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	slog.Info("interrupt signal received")
	signal.Stop(quit)
	cancel()
	slog.Info("attempting graceful shutdown")
	wg.Wait()
}
