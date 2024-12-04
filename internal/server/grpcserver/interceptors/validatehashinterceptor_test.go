package interceptors

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	pb "github.com/fasdalf/train-go-musthave-metrics/internal/common/proto/metrics"
)

func TestNewValidateHashInterceptor(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "valid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{
						constants.HeaderHashSHA256: "746865206b6579e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					}),
				)},
			wantResp: 10,
			wantErr:  false,
		},
		{
			name: "invalid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{constants.HeaderHashSHA256: "mock"}),
				),
			},
			wantResp: nil,
			wantErr:  true,
		},
	}

	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}
	req := pb.UpdateMetricsRequest{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewValidateHashInterceptor("the key")
			gotResp, err := got(tt.args.ctx, &req, &usi, nh)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHashInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ValidateHashInterceptor() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_makeContextWithHashValidMark(t *testing.T) {
	ctx := makeContextWithHashValidMark(context.Background())
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		t.Errorf("makeContextWithHashValidMark() metadata not found")
	}
	got := md.Get(hashPresentKey)
	if len(got) != 1 {
		t.Errorf("makeContextWithHashValidMark() metadata key %q not found", hashPresentKey)
	}
}
