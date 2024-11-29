package handlers

import (
	"context"
	"crypto/rsa"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

type grpcPoster struct {
	address       string
	key           string
	encryptionKey *rsa.PublicKey
}

func (p *grpcPoster) Post(ctx context.Context, idlog *slog.Logger, metrics []*apimodels.Metrics) error {
	mSlice := make([]*pb.Metric, len(metrics))
	for i, m := range metrics {
		mv := &pb.Metric{Id: m.ID}
		switch m.MType {
		case constants.CounterStr:
			mv.Type = pb.Metric_COUNTER
		case constants.GaugeStr:
			mv.Type = pb.Metric_GAUGE
		}
		if m.Delta != nil {
			mv.Delta = *m.Delta
		}
		if m.Value != nil {
			mv.Value = *m.Value
		}
		mSlice[i] = mv
	}
	mu := pb.UpdateMetricsRequest{Metrics: mSlice}
	// TODO: ##@@ add new metadata to ctx

	client, closeClient := p.newClient()
	defer closeClient()
	if _, err := client.UpdateMetrics(ctx, &mu); err != nil {
		idlog.Error("failed to invoke grpc", "error", err)
		return err
	}
	return nil
}

func NewGRPCPoster(address string, key string, encryptionKey *rsa.PublicKey) *grpcPoster {
	return &grpcPoster{
		address:       address,
		key:           key,
		encryptionKey: encryptionKey,
	}
}

func (p *grpcPoster) newClient() (pb.MetricsClient, func() error) {
	// It never returns error here.
	conn, _ := grpc.NewClient(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// TODO: ##@@ add new comporessors to conn
	// IRL use built in TLS transport and no compression
	client := pb.NewMetricsClient(conn)
	return client, conn.Close
}
