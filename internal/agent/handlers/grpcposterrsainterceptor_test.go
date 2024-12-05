package handlers

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/rsacrypt"
)

func Test_newGrpcPosterRSAInterceptor(t *testing.T) {
	type args struct {
		req any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				req: &pb.Metric{Id: "mock"},
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				req: "mock",
			},
			wantErr: true,
		},
	}
	ctx := context.Background()
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}
	_, pub := rsacrypt.GenerateKeyPair(2048)
	interceptor := newGrpcPosterRSAInterceptor(pub)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := interceptor(ctx, "any", tt.args.req, nil, nil, invoker)
			if (err != nil) != tt.wantErr {
				t.Errorf("newGrpcPosterRSAInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
