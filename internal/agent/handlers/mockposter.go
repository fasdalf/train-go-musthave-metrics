package handlers

import (
	"bytes"
	"context"
	"log/slog"
)

type mockPoster struct {
	Attempts int
	Cancel   context.CancelFunc
}

func (p *mockPoster) Post(ctx context.Context, idlog *slog.Logger, body *bytes.Buffer, key string, address string) error {
	p.Attempts--
	if (p.Attempts) <= 0 {
		p.Cancel()
	}
	return nil
}

func NewMockPoster(attempts int, cancel context.CancelFunc) *mockPoster {
	return &mockPoster{Attempts: attempts, Cancel: cancel}
}
