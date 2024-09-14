package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/cryptofacade"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"

	resty "github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

const URLTemplate = "http://%s/updates/"

type Retryer interface {
	Try(ctx context.Context, do func() error, isRetryable func(err error) bool) (int, error)
}

var ErrTransport = errors.New("resty error")

// SendMetrics sends pre collected metrics to server
func SendMetrics(ctx context.Context, s Storage, address string, r Retryer, key string, rateLimit int) {
	slog.Info("Sending metricUpdates")
	address = fmt.Sprintf(URLTemplate, address)

	w := newWorker(address, key)
	if _, err := r.Try(ctx, newProducer(ctx, s, w, rateLimit), isRecoverable); err != nil {
		slog.Info("SendMetrics error", "error", err)
	}
}

func isRecoverable(err error) bool {
	var netErr net.Error = nil
	if errors.As(err, &netErr) && netErr != nil {
		return true
	}
	return false
}

type workerFunc = func(ctx context.Context, id int, mCh <-chan *apimodels.Metrics) error

func newWorker(address string, key string) workerFunc {
	return func(ctx context.Context, id int, mCh <-chan *apimodels.Metrics) error {
		idlog := slog.With("workerFunc", "metricUpdates", "id", id)
		metricUpdates := make([]*apimodels.Metrics, 0)
		for m := range mCh {
			select {
			case <-ctx.Done():
				idlog.Error("context ended", "error", ctx.Err())
				return ctx.Err()
			default:
			}
			metricUpdates = append(metricUpdates, m)
		}

		idlog.Info("recieved metricUpdates", "count", len(metricUpdates))
		if len(metricUpdates) == 0 {
			return nil
		}

		client := resty.New()
		content, err := json.Marshal(metricUpdates)
		if err != nil {
			idlog.Error("marshal metricUpdates error", "error", err)
			return errors.Join(fmt.Errorf("encoding request: %w", err), ErrTransport)
		}
		body := new(bytes.Buffer)
		zb := gzip.NewWriter(body)
		_, err = zb.Write(content)
		if err != nil {
			idlog.Error("gzip writer error", "error", err)
			return fmt.Errorf("compressing request: %w", err)
		}
		_ = zb.Close()

		req := client.R()
		req.SetContext(ctx)
		req.SetHeader("Content-Encoding", "gzip")
		req.SetHeader("Accept-Encoding", "gzip")
		req.SetHeader("Content-Type", "application/json")

		if key != "" {
			hash := cryptofacade.Hash(body.Bytes(), []byte(key))
			req.SetHeader(constants.HashSHA256, hash)
		}

		req.SetBody(body)
		resp, err := req.Post(address)
		if err != nil {
			idlog.Error("send request error", "error", err)
			return fmt.Errorf("sending metrics: %w", err)
		}

		if resp != nil && resp.RawResponse.Body != nil {
			_ = resp.RawResponse.Body.Close()
		}

		if resp != nil && resp.IsError() {
			idlog.Error("response error", "error", resp.Error())
			return fmt.Errorf("sending metrics https status error: %d", resp.StatusCode())
		}
		idlog.Info("sent metricUpdates", "count", len(metricUpdates))
		return nil
	}
}

func newProducer(ctx context.Context, s Storage, w workerFunc, l int) func() error {
	return func() error {
		ch := make(chan *apimodels.Metrics)
		eg, innerCtx := errgroup.WithContext(ctx)
		for i := 0; i < l; i++ {
			eg.Go(func() error { return w(innerCtx, i, ch) })
		}

		for _, key := range s.ListCounters() {
			counter := int64(s.GetCounter(key))
			ch <- &apimodels.Metrics{
				ID:    key,
				MType: constants.CounterStr,
				Delta: &counter,
				Value: nil,
			}
		}
		for _, key := range s.ListGauges() {
			gauge := s.GetGauge(key)
			ch <- &apimodels.Metrics{
				ID:    key,
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
	address string,
	sendInterval time.Duration,
	retryer handlers.Retryer,
	key string,
	rateLimit int,
) {
	cb := func() {
		SendMetrics(ctx, storage, address, retryer, key, rateLimit)
		slog.Info(`sender sleeping`, `delay`, sendInterval)
	}
	loop(cb, ctx, wg, sendInterval)
}
