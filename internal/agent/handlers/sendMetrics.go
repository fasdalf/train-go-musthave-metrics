package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"

	"golang.org/x/sync/errgroup"
)

type Retryer interface {
	Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (int, error)
}

type MetricsPoster interface {
	Post(ctx context.Context, idlog *slog.Logger, metrics []*apimodels.Metrics) error
}

var ErrTransport = errors.New("resty error")

// SendMetrics sends pre collected metrics to server
func SendMetrics(ctx context.Context, s Storage, r Retryer, poster MetricsPoster, rateLimit int) {
	slog.Info("Sending metricUpdates")

	w := newWorker(poster)
	p := newProducer(ctx, s, w, rateLimit)
	if _, err := r.Try(ctx, p, isRecoverable); err != nil {
		slog.Info("SendMetrics error", "error", err)
	}
}

func isRecoverable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr != nil {
		return true
	}
	return false
}

type workerFunc = func(ctx context.Context, id int, mCh <-chan *apimodels.Metrics) error

func newWorker(poster MetricsPoster) workerFunc {
	return func(ctx context.Context, id int, mCh <-chan *apimodels.Metrics) error {
		idlog := slog.With("workerFunc", "metricUpdates", "id", id)
		metricUpdates := make([]*apimodels.Metrics, 0)
		for m := range mCh {
			metricUpdates = append(metricUpdates, m)
		}

		idlog.Info("received metricUpdates", "count", len(metricUpdates))
		if len(metricUpdates) == 0 {
			return nil
		}

		return poster.Post(ctx, idlog, metricUpdates)
	}
}

func newProducer(ctx context.Context, s Storage, w workerFunc, l int) func() error {
	return func() error {
		ch := make(chan *apimodels.Metrics)
		eg, innerCtx := errgroup.WithContext(ctx)
		for i := 0; i < l; i++ {
			eg.Go(func() error { return w(innerCtx, i, ch) })
		}

		for _, counterName := range s.ListCounters() {
			counter := int64(s.GetCounter(counterName))
			ch <- &apimodels.Metrics{
				ID:    counterName,
				MType: constants.CounterStr,
				Delta: &counter,
				Value: nil,
			}
		}
		for _, gaugeName := range s.ListGauges() {
			gauge := s.GetGauge(gaugeName)
			ch <- &apimodels.Metrics{
				ID:    gaugeName,
				MType: constants.GaugeStr,
				Delta: nil,
				Value: &gauge,
			}
		}

		close(ch)
		return eg.Wait()
	}
}

func SendMetricsLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	storage Storage,
	sendInterval time.Duration,
	retryer handlers.Retryer,
	poster MetricsPoster,
	rateLimit int,
) {
	cb := func() {
		SendMetrics(ctx, storage, retryer, poster, rateLimit)
		slog.Info(`sender sleeping`, `delay`, sendInterval)
	}
	loop(cb, ctx, wg, sendInterval)
}
