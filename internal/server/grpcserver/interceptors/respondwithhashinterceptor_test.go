package interceptors

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

func TestNewRespondWithHashInterceptor(t *testing.T) {
	type args struct {
		key     string
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{hashPresentKey: "mock"}))
	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}
	m := pb.UpdateMetricsRequest{}
	e := status.Error(codes.Unimplemented, "msg")
	eh := func(ctx context.Context, req any) (any, error) {
		return nil, e
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				key:     "filled",
				ctx:     ctx,
				req:     &m,
				info:    &usi,
				handler: nh,
			},
			wantResp: 10,
			wantErr:  false,
		},
		{
			name: "error",
			args: args{
				key:     "filled",
				ctx:     ctx,
				req:     &m,
				info:    &usi,
				handler: eh,
			},
			wantResp: nil,
			wantErr:  true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRespondWithHashInterceptor(tt.args.key)
			gotResp, err := got(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("mockInterceptor() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
