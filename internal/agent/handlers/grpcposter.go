package handlers

import (
	"context"
	"crypto/rsa"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/localip"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/hasher"
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

	mdMap := map[string]string{
		constants.HeaderRealIP: localip.GetLocalIP().String(),
	}
	if p.key != "" {
		mdMap[constants.HeaderHashSHA256] = hasher.Hash(&mu, []byte(p.key))
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(mdMap))

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
	conn, _ := grpc.NewClient(
		p.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(newGrpcPosterRSAInterceptor(p.encryptionKey)),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	)
	// IRL use built in TLS transport and signature instead of raw bytes field
	client := pb.NewMetricsClient(conn)
	return client, conn.Close
}
