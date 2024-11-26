package handlers

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
)

func TestSendMetricsLoop_EndToEnd(t *testing.T) {
	ms := metricstorage.NewMemStorageMuted()
	ms.UpdateCounter("testCounter", 10)
	ms.UpdateGauge("testGauge", 10.01)
	wg := &sync.WaitGroup{}

	t.Run("should send metrics", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		poster := NewMockPoster(3, cancel, []error{nil, nil, nil})
		wg.Add(1)
		SendMetricsLoop(
			ctx,
			wg,
			ms,
			2*time.Millisecond,
			retryattempt.NewOneAttemptRetryer(),
			poster,
			10,
		)

		if poster.Attempts > 0 {
			t.Errorf("expected poster to have 0 attempts, got %d", poster.Attempts)
		}
	})
	t.Run("should handle errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		poster := NewMockPoster(10, cancel, []error{errors.New("0"), &net.AddrError{}, errors.New("2"), errors.New("4"), errors.New("5")})
		wg.Add(1)
		SendMetricsLoop(
			ctx,
			wg,
			ms,
			2*time.Millisecond,
			retryattempt.NewRetryer([]time.Duration{1 * time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond}),
			poster,
			10,
		)

		if poster.Attempts > 0 {
			t.Errorf("expected poster to have 0 attempts, got %d", poster.Attempts)
		}
	})
	t.Run("should handle timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		poster := NewMockPoster(10, cancel, []error{&net.AddrError{}, &net.AddrError{}, errors.New("2"), errors.New("4"), errors.New("5")})
		wg.Add(1)
		SendMetricsLoop(
			ctx,
			wg,
			ms,
			2*time.Millisecond,
			retryattempt.NewRetryer([]time.Duration{1 * time.Millisecond, 200 * time.Millisecond, 300 * time.Millisecond}),
			poster,
			1,
		)

		if poster.Attempts <= 0 {
			t.Errorf("expected poster to have unused attempts, got %d", poster.Attempts)
		}
	})
}
