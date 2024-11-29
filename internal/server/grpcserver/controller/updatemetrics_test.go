package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

func TestMetricsServer_UpdateMetrics(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *pb.UpdateMetricsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateMetricsResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				r: &pb.UpdateMetricsRequest{
					Metrics: []*pb.Metric{
						{
							Id:    "ctr",
							Type:  pb.Metric_COUNTER,
							Delta: 10,
						},
						{
							Id:    "gg",
							Type:  pb.Metric_GAUGE,
							Value: 10.01,
						},
					},
				},
			},
			want:    &pb.UpdateMetricsResponse{},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				r: &pb.UpdateMetricsRequest{
					Metrics: []*pb.Metric{
						{
							Id:   "oops",
							Type: pb.Metric_UNDEFINED,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsServer{
				UnimplementedMetricsServer: pb.UnimplementedMetricsServer{},
				Storage:                    metricstorage.NewSavableModelStorage(metricstorage.NewMemStorage()),
			}
			got, err := s.UpdateMetrics(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
