package handlers

import (
	"context"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"sync"
	"testing"
	"time"
)

func TestSendMetricsLoop_EndToEnd(t *testing.T) {
	ms := metricstorage.NewMemStorageMuted()
	ms.UpdateGauge("testCounter", 10)
	ms.UpdateGauge("testGauge", 10.01)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	poster := NewMockPoster(3, cancel)

	t.Run("should send metrics", func(t *testing.T) {
		wg.Add(1)
		SendMetricsLoop(
			ctx,
			wg,
			ms,
			"",
			2*time.Millisecond,
			retryattempt.NewOneAttemptRetryer(),
			poster,
			"",
			10,
		)

		if poster.Attempts > 0 {
			t.Errorf("expected poster to have 0 attempts, got %d", poster.Attempts)
		}
	})
}
