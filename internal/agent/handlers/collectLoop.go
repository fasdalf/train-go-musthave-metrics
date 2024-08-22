package handlers

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type callback = func(s Storage, collectInterval time.Duration) error

func loop(c func(), ctx context.Context, wg *sync.WaitGroup, i time.Duration) {
	defer wg.Done()
	timer := time.NewTimer(i + 1)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		c()
		timer.Reset(i)
	}
}

func Collect(c callback, ctx context.Context, wg *sync.WaitGroup, storage Storage, collectInterval time.Duration) {
	cb := func() {
		if err := c(storage, collectInterval); err != nil {
			slog.Error(`collector error`, `err`, err)
		}
		slog.Info(`collector sleeping`, `delay`, collectInterval)
	}
	loop(cb, ctx, wg, collectInterval)
}
