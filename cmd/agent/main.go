package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
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
	memStorage := metricstorage.NewMemStorageMuted()
	retryer := retryattempt.NewRetryer([]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())
	// TODO: ##@@ extract to come package and cover with tests
	var poster handlers.MetricsPoster
	switch strings.ToLower(cfg.Protocol) {
	case "grpc":
		poster = handlers.NewGRPCPoster(cfg.Addr, cfg.HashKey, cfg.RSAKey)
	default:
		poster = handlers.NewNetHTTPPoster(cfg.Addr, cfg.HashKey, cfg.RSAKey)
	}

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go handlers.SendMetricsLoop(ctx, wg, memStorage, sendInterval, retryer, poster, cfg.RateLimit)
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
