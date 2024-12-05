package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/localip"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/grpcserver"
)

func TestNewGRPCPoster_Created(t *testing.T) {
	p := NewGRPCPoster(
		"localhost:2010",
		"key",
		nil,
	)
	if p == nil {
		t.Error("NewGRPCPoster() returned nil")
	}
}

func Test_grpcPoster_PostError(t *testing.T) {
	p := &grpcPoster{
		address:       "ftp://_:66333/",
		key:           "nil",
		encryptionKey: nil,
	}

	metrics := []*apimodels.Metrics{
		{
			ID:    "",
			MType: "",
			Delta: nil,
			Value: nil,
		},
	}

	if err := p.Post(context.Background(), slog.Default(), metrics); err == nil {
		t.Errorf("empty error")
	}
}

func Test_grpcPoster_PostSuccessE2E(t *testing.T) {
	ms := metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage())
	rt := retryattempt.NewOneAttemptRetryer()
	gs := grpcserver.NewGrpcServer(ms, nil, rt, "key", nil, nil)
	defer gs.Stop()
	port, _ := localip.GetFreePort()
	addr := fmt.Sprintf("localhost:%d", port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}
	go func() {
		if err := gs.Serve(listen); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	p := NewGRPCPoster(addr, "key", nil)

	i64 := int64(10)
	f64 := float64(10.01)
	metrics := []*apimodels.Metrics{
		{
			ID:    "ctr",
			MType: constants.CounterStr,
			Delta: &i64,
		},
		{
			ID:    "gge",
			MType: constants.GaugeStr,
			Value: &f64,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := p.Post(ctx, slog.Default(), metrics); err != nil {
		t.Errorf("failed to post: %v", err)
	}
}
