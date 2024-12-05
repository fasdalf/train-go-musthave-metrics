package handlers

import (
	"context"
	"log/slog"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
)

type mockPoster struct {
	Attempts int
	Cancel   context.CancelFunc
	Results  []error
}

func (p *mockPoster) Post(ctx context.Context, idlog *slog.Logger, metrics []*apimodels.Metrics) (err error) {
	i := len(p.Results) - p.Attempts

	if i >= 0 && i < len(p.Results) {
		err = p.Results[i]
	}
	p.Attempts--
	if (p.Attempts) <= 0 {
		p.Cancel()
	}
	return err
}

func NewMockPoster(attempts int, cancel context.CancelFunc, results []error) *mockPoster {
	return &mockPoster{Attempts: attempts, Cancel: cancel, Results: results}
}
