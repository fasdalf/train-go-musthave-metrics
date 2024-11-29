package controller

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

// UpdateMetrics updates metrics
func (s *MetricsServer) UpdateMetrics(ctx context.Context, r *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	if r != nil && len(r.Metrics) > 0 {
		updates := make([]apimodels.Metrics, len(r.Metrics))
		for i, m := range r.Metrics {
			u := apimodels.Metrics{
				ID:    m.Id,
				Delta: &m.Delta,
				Value: &m.Value,
			}
			switch m.Type {
			case pb.Metric_COUNTER:
				u.MType = constants.CounterStr
			case pb.Metric_GAUGE:
				u.MType = constants.GaugeStr
			}
			updates[i] = u
		}
		// todo: ##@@ add retryer
		err := s.Storage.SaveCommonModels(ctx, updates)
		if err != nil {
			slog.Error("can't save metrics on grpc", "error", err)
			return nil, status.Errorf(codes.Internal, "some error, see logs")
		}
		slog.Info("updated metrics", "count", len(updates))
	}
	return &pb.UpdateMetricsResponse{}, nil
}
